import {FC, useState} from 'react';
import {Avatar, Button, Image, Link, Placeholder, Section, Snackbar, Text} from '@telegram-apps/telegram-ui';
import {Page} from '@/components/Page.tsx';
import {bem} from '@/css/bem';

import {
    initDataRaw as _initDataRaw,
    initDataState as _initDataState,
    openTelegramLink,
    useSignal
} from '@telegram-apps/sdk-react';

import AuthService from '@/api/services/telegramAuthService.ts';

import './TelegramLoginPage.css';
import telegramLogo from './telegram-logo.svg';
import telegramDuck from './duck-telegram.webp'
import {Icon20ErrorCircleOutline} from "@vkontakte/icons";
import {useNavigate} from "react-router-dom";

const [, e] = bem('telegram-auth-prompt');

export const TelegramLoginPage: FC = () => {
    const initDataRaw = useSignal(_initDataRaw);
    const initDataState = useSignal(_initDataState);
    const user = initDataState?.user;
    const [loading, setLoading] = useState(false)
    const [snackbar, setSnackbar] = useState(false);
    const [error, setError] = useState('');
    const navigate = useNavigate();

    const fullName = user ? `${user.first_name ?? ''} ${user.last_name ?? ''}`.trim() : '';

    const handleLogin = async () => {
        setError('')
        if (!initDataRaw) {
            setError('Не удалось получить данные для входа');
            setSnackbar(true)
            return
        }
        try {
            setLoading(true)
            const response = await AuthService.loginWithTelegram({init_data: initDataRaw});
            console.log('Успешный вход:', response);
            const {access_token, refresh_token} = response.data;
            console.log('tokens: ', access_token)
            sessionStorage.setItem('access_token', access_token || '');
            sessionStorage.setItem('refresh_token', refresh_token || '');
            navigate("/bank")
        } catch (err) {
            console.error('Ошибка входа:', err);
            setError('Ошибка входа. Попробуйте снова');
            setSnackbar(true)
        } finally {
            setLoading(false)
        }
    };

    return (
        <Page>
            {snackbar &&
                <Snackbar
                    onClose={() => setSnackbar(false)}
                    before={<Icon20ErrorCircleOutline/>}
                    description={error}
                    children={'Ошибка'}
                    style={{paddingBottom: "20px"}}

                />
            }
            <Placeholder>
            <Section className={e('container')}>
                <Placeholder
                    children={
                        user?.photo_url ? (
                            <Avatar size={96} src={user.photo_url}/>
                        ) : (
                            <Avatar size={96}/>
                        )
                    }
                    header={fullName || 'Telegram Login'}
                    description={
                        <Text>

                            {fullName
                                ? <Text>
                                    <Text>Вы входите как </Text>
                                    <Link onClick={() => {
                                        if (openTelegramLink.isAvailable()) {
                                            openTelegramLink(`https://t.me/${user?.username}`);
                                        }
                                    }}
                                    style={{cursor: "pointer"}}
                                    >
                                        @{user?.username}
                                    </Link>
                                </Text>
                                : 'Чтобы продолжить, необходимо авторизоваться в Telegram.'}
                        </Text>
                    }
                    action={
                        <Button
                            loading={loading}
                            size="l"
                            mode="filled"
                            stretched
                            className={e('button')}
                            onClick={handleLogin}
                        >
                            <Button
                                before={<Image size={28} src={telegramLogo}
                                               style={{backgroundColor: 'transparent', boxShadow: "none"}}/>}
                            >
                                {fullName ? `Продолжить как ${fullName}` : 'Войти в Telegram'}
                            </Button>


                        </Button>
                    }
                />
            </Section>
                <img
                    style={{
                        width: "200px",
                        height: "200px"
                    }}
                    alt={'Telegram duck'}
                    src={telegramDuck}
                />
            </Placeholder>

        </Page>
    );
};
