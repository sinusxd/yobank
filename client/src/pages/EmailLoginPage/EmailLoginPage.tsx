import { FC, useState, useEffect } from 'react';
import { Page } from '@/components/Page';
import { Button, Input, PinInput, Placeholder } from '@telegram-apps/telegram-ui';
import { EmailCodeRequest } from '@/api/models/request/emailCodeRequest';
import { VerifyCodeRequest } from '@/api/models/request/verifyCodeRequest';
import { VerifyCodeResponse } from '@/api/models/response/verifyCodeResponse';
import AuthService from '@/api/services/telegramAuthService.ts';
import styles from './EmailLoginPage.module.css';

export const EmailLoginPage: FC = () => {
    const [step, setStep] = useState<'email' | 'code'>('email');
    const [email, setEmail] = useState('');
    const [code, setCode] = useState<number[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [shake, setShake] = useState(false);

    const requestCode = async () => {
        setLoading(true);
        setError('');
        try {
            const request: EmailCodeRequest = { email };
            await AuthService.requestCode(request);
            setStep('code');
        } catch (err: any) {
            setError(err?.response?.data?.message || err.message || 'Ошибка при отправке кода');
        } finally {
            setLoading(false);
        }
    };

    const verifyCode = async () => {
        setLoading(true);
        setError('');
        try {
            const request: VerifyCodeRequest = { email, code: code.join('') };
            const response = await AuthService.verifyCode(request);
            const { access_token, refresh_token }: VerifyCodeResponse = response.data;

            sessionStorage.setItem('access_token', access_token || '');
            sessionStorage.setItem('refresh_token', refresh_token || '');
            // TODO: router.push('/') или callback
        } catch (err: any) {
            setError('Неверный или просроченный код');

            // 🎯 Вибрация + shake
            if (navigator.vibrate) {
                navigator.vibrate(200);
            }
            setShake(true);
            setTimeout(() => setShake(false), 400);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (code.length === 6) {
            verifyCode();
        }
    }, [code]);

    return (
        <Page back>
            <div style={{ padding: '1.5rem' }}>
                <Placeholder
                    header={step === 'email' ? 'Вход по почте' : 'Введите код'}
                    description={
                        step === 'email'
                            ? 'Мы отправим код подтверждения на вашу почту'
                            : `На ${email} отправлен код`
                    }
                />

                {step === 'email' && (
                    <div style={{ display: 'flex', flexDirection: 'column', gap: '2rem' }}>
                        <Input
                            placeholder="Введите вашу почту"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            disabled={loading}
                            type="email"
                        />
                        <Button
                            stretched
                            onClick={requestCode}
                            loading={loading}
                            disabled={!email}
                        >
                            Получить код
                        </Button>
                    </div>
                )}

                {step === 'code' && (
                    <div
                        className={shake ? styles.shake : ''}
                        style={{
                            width: '100%',
                            maxWidth: 240,
                            margin: '2rem auto 0',
                        }}
                    >
                        <PinInput
                            label="Введите код из письма"
                            pinCount={6}
                            value={code}
                            onChange={(val) => {
                                setCode(val);
                                if (error) setError('');
                            }}
                        />
                    </div>
                )}

                {error && (
                    <p style={{ color: 'red', marginTop: '1.5rem', textAlign: 'center' }}>
                        {error}
                    </p>
                )}
            </div>
        </Page>
    );
};
