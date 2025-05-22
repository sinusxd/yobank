import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { User } from "@/api/services/userService.ts";

// Попытка загрузить user из sessionStorage
const savedUser = sessionStorage.getItem("user");
const initialState: User | null = savedUser ? JSON.parse(savedUser) : null;

const userSlice = createSlice({
    name: 'user',
    initialState,
    reducers: {
        setUser: (_state, action: PayloadAction<User>) => {
            sessionStorage.setItem("user", JSON.stringify(action.payload));
            return action.payload;
        },
        clearUser: () => {
            sessionStorage.removeItem("user");
            return null;
        },
    },
});

export const { setUser, clearUser } = userSlice.actions;
export default userSlice.reducer;
