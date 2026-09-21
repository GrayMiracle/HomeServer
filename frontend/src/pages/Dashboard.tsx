// frontend/src/pages/Dashboard.tsx
import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { ArrowLeft } from 'lucide-react'
import { apiFetch } from '../api'
import { formatDate } from '../utils'

type Activity = {
    id: number
    eventType: string
    filePath: string
    userId: string
    deviceId: string
    ipAddress: string
    createdAt: string
}

const EVENT_LABELS: Record<string, string> = {
    login: 'Logged in',
    refresh: 'Session refreshed',
    logout: 'Logged out',
    pin_verified: 'PIN verified',
    pin_failed: 'PIN attempt failed',
    view: 'Viewed file',
    download: 'Downloaded file',
    upload: 'Uploaded file',
    delete: 'Deleted file',
    stream: 'Streamed file',
}

export default function Dashboard() {
    const navigate = useNavigate()
    
    const [activities, setActivities] = useState<Activity[]>([])
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState('')


    useEffect(() => {
        async function load() {
            try {
                const res = await apiFetch('/activity')
                if (!res.ok) throw new Error('Failed to load')
                const data = await res.json()
                setActivities(Array.isArray(data) ? data : [])
            } catch {
                setError('Failed to load activity')
            } finally {
                setLoading(false)
            }
        }
        load()
    }, [])

    return (
    <div className="min-h-screen bg-[#0a0a0a] text-white">
        {/* Top bar — matches FileExplorer */}
        <div className="relative flex items-center px-6 py-4 bg-[#111] border-b border-[#1f1f1f]">
            <button
                onClick={() => navigate(-1)}
                className="flex items-center gap-2 bg-[#1f1f1f] hover:bg-[#2a2a2a] text-white text-sm font-medium px-4 py-2 rounded-lg transition-colors duration-150"
            >
                <ArrowLeft size={14} />
            </button>
            <h1 className="absolute left-1/2 -translate-x-1/2 text-3xl font-bold">Activity</h1>
        </div>
        {/* Content */}
        <div className="p-6 flex flex-col gap-2">
            {loading && <p className="text-sm text-[#666]">Loading...</p>}
            {error && <p className="text-sm text-[#e03030]">{error}</p>}
            {!loading && !error && activities.length === 0 && (
                <p className="text-sm text-[#666]">No activity yet.</p>
            )}
            {activities.map(a => (
                <div
                    key={a.id}
                    className="flex items-center justify-between bg-[#111] border border-[#1f1f1f] rounded-lg px-4 py-3"
                >
                    <div className="text-left">
                        <p className="text-sm font-medium">{EVENT_LABELS[a.eventType] ?? a.eventType}</p>
                        {a.filePath && <p className="text-xs text-[#666] mt-1">{a.filePath}</p>}
                    </div>
                    <div className="text-right">
                        <p className="text-xs text-[#666]">{formatDate(a.createdAt)}</p>
                        <p className="text-xs text-[#555] mt-1">{a.ipAddress}</p>
                    </div>
                </div>
            ))}
        </div>
    </div>
)
}