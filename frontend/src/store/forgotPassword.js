import { defineStore } from "pinia";
import axios from 'axios';

const BASE_URL = `${import.meta.env.VITE_API_BASE_URL}/auth/forgot-password`;

const START_FORGOT_PWD_API = BASE_URL + '/start';
const VERIFY_OTP_API = BASE_URL + '/verify';
const RESET_PASSWORD_API = BASE_URL + '/reset';
const RESEND_OTP_API = BASE_URL + '/resend';

export const useForgotPasswordStore = defineStore('forgotPassword', {
    state: () => ({
        flowId: localStorage.getItem('flowId') || null,
        email: localStorage.getItem('resetEmail') || '',
    }),

    actions: {
        setFlow(flowId, email) {
            this.flowId = flowId;
            this.email = email;
            localStorage.setItem('flowId', flowId);
            localStorage.setItem('resetEmail', email);
        },

        resetFlow() {
            this.flowId = null;
            this.email = null;
            localStorage.removeItem('flowId');
            localStorage.removeItem('resetEmail');
        },

        async startForgotPassword(email) {
            try {
                const res = await axios.post(START_FORGOT_PWD_API, { email });
                if(res.status === 201) {
                    this.setFlow(res.data.flowId, email);
                    return { success: true };
                }
            } catch (error) {
                if(error.response?.status === 429) {
                    this.resetFlow();
                    throw { redirect: true, message: 'Too many requests' }
                }
                throw error;
            }
        },

        async verifyOtp(otp) {
            try {
                const payload = { flowId: this.flowId, email: this.email, otp };
                const res = await axios.post(VERIFY_OTP_API, payload);
                if(res.status === 200) {
                    return { success: true };
                }
            } catch (error) {
                if(error.response?.status === 429) {
                    this.resetFlow();
                    throw { redirect: true, message: 'Too many requests' }
                }
                throw error;
            }
        },

        async resendOtp() {
            try {
                const payload = { flowId: this.flowId, email: this.email };
                const res = await axios.post(RESEND_OTP_API, payload);
                return res.data;
            } catch (error) {
                if(error.response?.status === 429) {
                    this.resetFlow();
                    throw { redirect: true, message: 'Too many requests' }
                }
                throw error;
            }
        },

        async resetPassword(newPassword) {
            try {
                const payload = { flowId: this.flowId, email: this.email, newPassword };
                const res = await axios.patch(RESET_PASSWORD_API, payload);
                if(res.status === 200) {
                    this.resetFlow();
                    return { success: true };
                }
            } catch (error) {
                if(error.response?.status === 429) {
                    this.resetFlow();
                    throw { redirect: true, message: 'Too many requests' }
                }
                throw error;
            }
        }
    }
})