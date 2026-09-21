import { useState } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import { API } from '../config'

// Get or Create Device ID by checking localStorage
function getDeviceId() {
    let deviceId = localStorage.getItem('deviceId')
    if (!deviceId) {
        deviceId = generateUUID()
        localStorage.setItem('deviceId', deviceId)
    }
    return deviceId
}

function generateUUID(): string {
    if (typeof crypto !== 'undefined' && crypto.randomUUID) {
        return crypto.randomUUID()
    }
    // RARE fallback when crypto.randomUUID is unavailable
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, c => {
        const r = (Math.random() * 16) | 0
        const v = c === 'x' ? r : (r & 0x3) | 0x8
        return v.toString(16)
    })
}

// Login function
export default function Login() {
    const [username, setUsername] = useState('')
    const [password, setPassword] = useState('')
    const [error, setError] = useState('')
    const [loading, setLoading] = useState(false)
    const navigate = useNavigate()
    const location = useLocation()

    const from = (location.state as { from?: Location })?.from?.pathname || '/'

    async function handleLogin(e: React.FormEvent) {
        e.preventDefault()
        setError('')
        setLoading(true)

        try {
            const res = await fetch(`${API}/auth/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({
                    username,
                    password,
                    deviceId: getDeviceId(),
                    deviceName: navigator.userAgent,
                }),
            })

            if (!res.ok) {
                const message = await res.text()
                setError(message || 'Login failed')
                return
            }

            const data = await res.json()
            localStorage.setItem('refreshToken', data.refreshToken)
            localStorage.setItem('userId', String(data.userId))
            // Finally navigate to where user came from
            navigate(from, { replace: true })

        } catch (err) {
            console.log('Login error:', err)
            setError('Failed to reach server')
        } finally {
            setLoading(false)
        }
    }

    return (
        <div className="min-h-screen flex items-center justify-center">
            <form onSubmit={handleLogin} className="flex flex-col gap-3 w-64">
                <h1 className="text-xl mb-2">HomeCloud Login</h1>

                {error && <p className="text-red-500 text-sm">{error}</p>}

                <input
                    type="text"
                    placeholder="Username"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    className="bg-gray-800 px-3 py-2 rounded"
                    required
                />
                <input
                    type="password"
                    placeholder="Password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className="bg-gray-800 px-3 py-2 rounded"
                    required
                />

                <button
                    type="submit"
                    disabled={loading}
                    className="bg-red-600 px-3 py-2 rounded disabled:opacity-50"
                >
                {loading ? 'Logging in...' : 'Log in'}
                </button>
            </form>
        </div>
    )
}