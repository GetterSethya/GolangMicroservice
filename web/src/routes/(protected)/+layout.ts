import { goto } from '$app/navigation';
import { UserRepository } from '@/repository/user';
import type { LayoutLoad } from './$types';
import { readable } from 'svelte/store';

export const load: LayoutLoad = async () => {
	const accessToken = localStorage.getItem('accessToken');
	if (!accessToken) {
		goto('/login');
	}

	const userRepo = new UserRepository();
	const {res} = await userRepo.getLocalUserData();
    const localUser =res.data?.user
    if (!localUser) {
        localStorage.removeItem("accessToken")
        localStorage.removeItem("refreshToken")
        goto('/login')
    }
    return readable({
        user:localUser
    })
};
