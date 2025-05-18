import {type FC} from 'react';

import {Cell, Image, List, Placeholder, Section, Title} from '@telegram-apps/telegram-ui';
import {Page} from '@/components/Page.tsx';

import bankLogo from './bank.gif';
import mailLogo from './mail.png';
import telegramLogo from "./telegram.png"
import {Link} from "@/components/Link/Link.tsx";

export const LoginPage: FC = () => {

    return (
        <Page back={false}>

            <div style={{display: 'flex', flexDirection: 'column', alignItems: 'center', padding: '2rem 1rem 1rem'}}>
                <Placeholder
                    header={
                        <Title
                            level="1"
                            weight="1"
                        >
                           Добро пожаловать в YoBank!
                        </Title>
                    }
                >
                    <img
                        alt="Telegram sticker"
                        src={bankLogo}
                        style={{width: 120, height: 120, borderRadius: 20, marginBottom: '1rem'}}
                    />
                </Placeholder>
            </div>

            <List>
                <Section header="Выберите способ входа">
                    <Link to={"/telegram-login"}>
                        <Cell
                            before={<Image src={telegramLogo} style={{backgroundColor: '#24A1DE'}}/>}
                            subtitle="Используется ваш telegram-аккаунт"
                        >
                            Продолжить с Telegram
                        </Cell>
                    </Link>
                    <Link to={"/email-login"}>
                        <Cell
                            before={<Image src={mailLogo} style={{backgroundColor: '#24A1DE'}}/>}
                            subtitle="Используется ваша почта"
                        >
                            Продолжить с Email
                        </Cell>
                    </Link>
                </Section>
            </List>
        </Page>
    );
};
