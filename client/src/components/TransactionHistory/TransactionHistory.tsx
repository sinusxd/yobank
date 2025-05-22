import { FC, useEffect, useRef, useState } from "react";
import { Avatar, Badge, Cell, List, Section, Skeleton, Text, Title } from "@telegram-apps/telegram-ui";
import WalletService from "@/api/services/walletService.ts";
import TransferService, { Transfer } from "@/api/services/transferService.ts";
import dayjs from "dayjs";
import { useSelector } from "react-redux";
import { RootState } from "@/store/store.ts";
import { User } from "@/api/services/userService.ts";

interface UserInfo {
    username: string;
    avatarUrl: string;
}

export const TransactionHistory: FC = () => {
    const currentUser = useSelector<RootState, User | null>((state) => state.user);
    const [transfers, setTransfers] = useState<Transfer[]>([]);
    const [userWalletIds, setUserWalletIds] = useState<number[]>([]);
    const [userInfoMap, setUserInfoMap] = useState<Record<number, UserInfo>>({});
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const fallbackAvatarsRef = useRef<Map<number, string>>(new Map());

    const getStableFallbackAvatar = (walletId: number): string => {
        const map = fallbackAvatarsRef.current;
        if (!map.has(walletId)) {
            const seed = walletId % 1000000;
            map.set(walletId, `https://avatars.githubusercontent.com/u/${seed}?v=4`);
        }
        return map.get(walletId)!;
    };

    useEffect(() => {
        const loadHistory = async () => {
            if (!currentUser) {
                setError("Пользователь не найден");
                return;
            }

            setLoading(true);
            try {
                const { data: wallets } = await WalletService.getUserWallets();
                if (wallets.length === 0) {
                    setError("Кошельки пользователя не найдены");
                    return;
                }

                const walletIds = wallets.map(w => w.id);
                setUserWalletIds(walletIds);

                const allTransfers: Transfer[] = [];

                for (const id of walletIds) {
                    const transfers = await TransferService.getTransferHistory(id);
                    allTransfers.push(...transfers);
                }

                allTransfers.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime());
                setTransfers(allTransfers);

                const participantIds = new Set<number>();
                allTransfers.forEach(tx => {
                    participantIds.add(tx.senderWalletId);
                    participantIds.add(tx.receiverWalletId);
                });

                const resolved: Record<number, UserInfo> = { ...userInfoMap };

                await Promise.all(
                    Array.from(participantIds).map(async (id) => {
                        if (resolved[id]) return;

                        try {
                            const { username, avatarUrl } = await TransferService.getReceiverUsername(id);
                            resolved[id] = {
                                username,
                                avatarUrl: avatarUrl || getStableFallbackAvatar(id),
                            };
                        } catch {
                            resolved[id] = {
                                username: "неизвестно",
                                avatarUrl: getStableFallbackAvatar(id),
                            };
                        }
                    })
                );

                setUserInfoMap(resolved);
            } catch {
                setError("Ошибка загрузки истории переводов");
            } finally {
                setLoading(false);
            }
        };

        loadHistory();
    }, [currentUser]);

    if (error) {
        return (
            <Section>
                <Cell style={{ color: "var(--tgui--destructive_text_color)" }}>{error}</Cell>
            </Section>
        );
    }

    return (
        <List>
            <Section header="История переводов">
                <Skeleton visible={loading}>
                    {transfers.length > 0 ? (
                        transfers.map((transfer) => {
                            const isIncoming = userWalletIds.includes(transfer.receiverWalletId);
                            const counterpartyId = isIncoming ? transfer.senderWalletId : transfer.receiverWalletId;
                            const user = userInfoMap[counterpartyId];
                            const username = user?.username ?? "неизвестно";
                            const avatarUrl = user?.avatarUrl ?? getStableFallbackAvatar(counterpartyId);

                            const sign = isIncoming ? "+" : "-";
                            const direction = isIncoming ? `от ${username}` : `${username}`;
                            const subtitle = isIncoming ? `Получено ${direction}` : `Отправлено ${direction}`;

                            const amountStyle = {
                                color: isIncoming ? "var(--tgui--accent_text_color)" : "var(--tgui--destructive_text_color)",
                                fontWeight: 600,
                            };

                            return (
                                <Cell
                                    key={transfer.id}
                                    before={<Avatar size={48} src={avatarUrl} />}
                                    titleBadge={<Badge type="dot" />}
                                    subhead={dayjs(transfer.createdAt).format("DD.MM.YYYY HH:mm")}
                                    subtitle={subtitle}
                                    description={
                                        <span style={amountStyle}>
                                            {sign}{(transfer.amount / 100).toFixed(2)} {transfer.currency}
                                        </span>
                                    }
                                >
                                    <Title level="3">Перевод</Title>
                                </Cell>
                            );
                        })
                    ) : (
                        <Text>Нет переводов</Text>
                    )}
                </Skeleton>
            </Section>
        </List>
    );
};
