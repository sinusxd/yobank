import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import {User} from "@/api/services/userService.ts";



// Главное: явно указать тип состояния в createSlice
const userSlice = createSlice({
    name: 'user',
    initialState: null as User | null,
    reducers: {
        setUser: (_state, action: PayloadAction<User>) => action.payload,
        clearUser: () => null,
    },
});

export const { setUser, clearUser } = userSlice.actions;
export default userSlice.reducer;
