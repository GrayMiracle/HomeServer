import { useState, useEffect } from 'react'
import type { ReactNode } from 'react'
import { Navigate, useLocation } from 'react-router-dom'
import { apiFetch } from '../api'

type Props = {
    children: ReactNode
}

export default function AuthGate({ children }: Props) {
    const location = useLocation()
    const [loading, setLoading] = useState(true)
    const [authenticated, setAuthenticated] = useState(false)

    useEffect(() => {
        if (localStorage.getItem('refreshToken') === null) {
            setLoading(false)
            setAuthenticated(false)
            return
        }

        async function checkAuth() {
            try {
                const res = await apiFetch('/auth/me')
                setAuthenticated(res.ok)
            } catch {
                setAuthenticated(false)
            } finally {
                setLoading(false)
            }
        }

        checkAuth()
    }, [])

    if (loading) {
        return (
            <div className="min-h-screen flex items-center justify-center bg-[#0a0a0a] text-white">
                Loading...
            </div>
        )
    }
    
    if (!authenticated) return <Navigate to="/login" replace state={{ from: location }} />
    return <>{children}</>
}