import type { ReactNode } from 'react'
import { Navigate } from 'react-router-dom'

export default function GuestGate({ children }: { children: ReactNode }) {
    const loggedIn = localStorage.getItem('refreshToken') !== null

    if (loggedIn) {
        return <Navigate to="/" replace />
    }

    return <>{children}</>
}