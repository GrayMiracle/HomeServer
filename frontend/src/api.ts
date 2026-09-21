// Auth handling for every request

import { API } from './config'

// Get stored access token, if valid return and use. If not try to refresh
export async function apiFetch(input: string, init: RequestInit = {}): Promise<Response> {
    // Add credentials for cookies
    const authInit: RequestInit = {
        ...init,
        credentials: 'include', // Includes HttpOnly access token cookie automatically
    }

    const res = await fetch(`${API}${input}`, authInit)

    // If access token is valid, return response
    if (res.status !== 401) return res

    // Try to refresh token if we get 401
    const refreshed = await tryRefresh()
    if (!refreshed) {
        logout()
        return res
    }

    return fetch(`${API}${input}`, authInit)
}

// Refresh function
async function tryRefresh(): Promise<boolean> {
    const refreshToken = localStorage.getItem('refreshToken')
    const userId = localStorage.getItem('userId')
    const deviceId = localStorage.getItem('deviceId')

    if (!refreshToken || !userId || !deviceId) return false

    try {
        const res = await fetch(`${API}/auth/refresh`, {
            method: 'POST',
            credentials: 'include',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                refreshToken,
                userId: Number(userId),
                deviceId,
            })
        })
    
        // res.ok is true if worked, false if refresh token is invalid/expired
        return res.ok
    } catch {
        return false
    }
}

export async function logout() {
    try {
        await fetch(`${API}/auth/logout`, {
            method: 'POST',
            credentials: 'include',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                userId: Number(localStorage.getItem('userId')),
                deviceId: localStorage.getItem('deviceId'),
            })
        })
    } catch {}
    localStorage.removeItem('refreshToken')
    localStorage.removeItem('userId')
    localStorage.removeItem('deviceId')
    window.location.href = '/login'
}