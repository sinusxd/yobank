import { FC, useState, useEffect, useRef } from 'react';
import { Section, Placeholder, Button, Input, PinInput } from '@telegram-apps/telegram-ui';
import { bem } from '@/css/bem';
import AuthService from '@/api/services/telegramAuthService';
import './EmailLoginPage.css';
import {EmailCodeRequest} from "@/api/models/request/emailCodeRequest.ts";
import {VerifyCodeResponse} from "@/api/models/response/verifyCodeResponse.ts";
import {VerifyCodeRequest} from "@/api/models/request/verifyCodeRequest.ts";
import {Page} from "@/components/Page.tsx";

const [, e] = bem('email-login-page');

export const EmailLoginPage: FC = () => {
    const [step, setStep] = useState<'email' | 'code'>('email');
    const [email, setEmail] = useState('');
    const [code, setCode] = useState<number[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [shake, setShake] = useState(false);
    const first = useRef(true);

    const requestCode = async () => {
        setLoading(true);
        setError('');
        try {
            await AuthService.requestCode({ email } as EmailCodeRequest);
            setStep('code');
        } catch (e: any) {
            setError(e?.response?.data?.message || e.message || 'Ошибка отправки кода');
        } finally {
            setLoading(false);
        }
    };

    const verifyCode = async () => {
        setLoading(true);
        setError('');
        try {
            const { data } = await AuthService.verifyCode({ email, code: code.join('') } as VerifyCodeRequest);
            const { access_token, refresh_token }: VerifyCodeResponse = data;
            sessionStorage.setItem('access_token', access_token || '');
            sessionStorage.setItem('refresh_token', refresh_token || '');
        } catch {
            setError('Неверный или просроченный код');
            setShake(true);
            if ('vibrate' in navigator) navigator.vibrate(200);
            setTimeout(() => {
                setShake(false);
                setCode([]);
            }, 600);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (code.length === 6) verifyCode();
    }, [code]);

    return (
        <Page back>
            {step === 'email' && (
                <Section>
                    <Placeholder
                        className={e('placeholder')}
                        header="Вход по почте"
                        description="Мы отправим код подтверждения на вашу почту"
                    />
                    <Input
                        className={e('input')}
                        type="email"
                        placeholder="Введите вашу почту"
                        value={email}
                        onChange={v => setEmail(v.target.value)}
                        disabled={loading}
                    />
                    <Button
                        className={e('button')}
                        stretched
                        loading={loading}
                        disabled={!email}
                        onClick={requestCode}
                    >
                        Получить код
                    </Button>
                    {error && <p className={e('error')}>{error}</p>}
                </Section>
            )}

            {step === 'code' && (
                <Section className={e('pin-section')}>
                    <Placeholder
                        header="Введите код"
                        description={`На ${email} отправлен код`}
                    />
                    <PinInput
                        className={e('pin-input', { shake })}
                        pinCount={6}
                        value={code}
                        onChange={vals => {
                            if (first.current) {
                                first.current = false;
                                return;
                            }
                            setCode(vals);
                            if (error) setError('');
                        }}
                    />
                    {error && <p className={e('error')}>{error}</p>}
                </Section>
            )}
        </Page>
    );
};
