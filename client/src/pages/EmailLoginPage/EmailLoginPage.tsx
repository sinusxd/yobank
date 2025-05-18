import {FC, useEffect, useRef, useState} from 'react';
import {Button, Input, List, PinInput, Placeholder, Section, Snackbar} from '@telegram-apps/telegram-ui';
import {bem} from '@/css/bem';
import AuthService from '@/api/services/telegramAuthService';
import './EmailLoginPage.css';
import {EmailCodeRequest} from "@/api/models/request/emailCodeRequest.ts";
import {VerifyCodeRequest} from "@/api/models/request/verifyCodeRequest.ts";
import {VerifyCodeResponse} from "@/api/models/response/verifyCodeResponse.ts";
import {Page} from "@/components/Page.tsx";
import {retrieveLaunchParams} from "@telegram-apps/sdk-react";
import codeLogo from './code.gif'
import {Icon20ErrorCircleOutline} from "@vkontakte/icons";
import {isValidEmail} from "@/utils/validateEmail.ts";

const [, e] = bem('email-login-page');

export const EmailLoginPage: FC = () => {
    const [step, setStep] = useState<'email' | 'code'>('email');
    const [email, setEmail] = useState('');
    const [code, setCode] = useState<number[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [shake, setShake] = useState(false);
    const [animating, setAnimating] = useState(false)
    const [pinKey, setPinKey] = useState(0);
    const skipFirst = useRef(true);
    const [snackbar, setSnackbar] = useState(false);
    const [inputStatus, setInputStatus] = useState<'default' | 'error' | 'focused'>('default')

    const launchParams = retrieveLaunchParams();
    const {tgWebAppPlatform: platform} = launchParams;

    useEffect(() => {
        if (platform === 'ios') {
            document.body.classList.add('ios');
        } else if (platform === 'macos') {
            document.body.classList.add('macos');
        }
    }, [platform]);

    const requestCode = async () => {
        setLoading(true);
        setError('');
        try {
            if (!isValidEmail(email)) {
                setError('Введите корректный email');
                setInputStatus('error')
                setShake(true)
                setTimeout(() => {
                    setShake(false)
                }, 1000)
                setTimeout(() => setInputStatus('default'), 500)
                setSnackbar(true)
                return;
            }
            await AuthService.requestCode({email} as EmailCodeRequest);
            setStep('code');
        }
        catch {
            setError('Не удалось отправить код на почту');
            setInputStatus('error')
            setShake(true)
            setTimeout(() => {
                setShake(false)
            }, 1000)
            setTimeout(() => setInputStatus('default'), 500)
            setSnackbar(true)
        }
        finally {
            setLoading(false);
        }
    };

    const verifyCode = async () => {
        setLoading(true);
        setAnimating(true)
        setError('');
        try {
            const {data} = await AuthService.verifyCode(
                {email, code: code.join('')} as VerifyCodeRequest,
            );
            const {access_token, refresh_token}: VerifyCodeResponse = data;
            sessionStorage.setItem('access_token', access_token || '');
            sessionStorage.setItem('refresh_token', refresh_token || '');
        } catch {
            setError('Неверный или просроченный код');
            setAnimating(false)
            setShake(true);
            if ('vibrate' in navigator) navigator.vibrate(200);
            setTimeout(() => {
                setShake(false);
                skipFirst.current = true;
                setCode([]);
                setPinKey(k => k + 1);
            }, 600);
        } finally {
            setLoading(false);
            setAnimating(false)
        }
    };

    useEffect(() => {
        if (code.length === 6) verifyCode();
    }, [code]);

    return (
        <Page back>
            {snackbar &&
                <Snackbar
                    onClose={() => setSnackbar(false)}
                    before={<Icon20ErrorCircleOutline/>}
                    description={error}
                    children={'Ошибка'}
                    style={{paddingBottom: "20px"}}

                />
            }
            {step === 'email' && (
                <div style={{padding: "10px"}}>
                    <Placeholder
                        className={e('placeholder')}
                        header="Вход по почте"
                        description="Мы отправим код подтверждения на вашу почту"
                    />
                    <>
                        <List style={{
                            padding: "35px",
                            display: "flex",
                            flexDirection: "column",
                            alignItems: "center",
                            gap: "10px"
                        }}
                              className={e('list', {shake})}
                        >
                            <div style={{width: "100%"}}>
                                <Input
                                    status={inputStatus || 'default'}
                                    style={{textAlign: 'center'}}
                                    type="email"
                                    placeholder="Введите вашу почту"
                                    value={email}
                                    onChange={v => setEmail(v.target.value)}
                                    disabled={loading}
                                />
                                <Button
                                    style={{marginTop: "20px"}}
                                    className={e('button')}
                                    stretched
                                    loading={loading}
                                    disabled={!email}
                                    onClick={requestCode}
                                >
                                    Получить код
                                </Button>
                            </div>
                            <img
                                width={"200"}
                                height={"200"}
                                alt={"duck"}
                                src={codeLogo}
                            />
                        </List>
                    </>
                </div>
            )}

            {step === 'code' && (
                <Section>
                    <PinInput
                        label="Введите код"
                        key={pinKey}
                        className={e('pin-input', {shake, animating})}
                        pinCount={6}
                        value={code}
                        onChange={vals => {
                            if (skipFirst.current) {
                                skipFirst.current = false;
                                return;
                            }
                            setCode(vals as number[]);
                            if (error) setError('');
                        }}
                    />
                    {error && step === 'code' && <p className={e('error')}>{error}</p>}
                </Section>
            )}
        </Page>
    );
};
