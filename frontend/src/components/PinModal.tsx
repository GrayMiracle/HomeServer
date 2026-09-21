import { useState } from 'react'
import { apiFetch } from '../api'

type Props = {
    onSuccess: () => void
    onCancel: () => void
}

export default function PinModal({ onSuccess, onCancel }: Props) {
    const [pin, setPin] = useState('')
    const [error, setError] = useState('')
    const [loading, setLoading] = useState(false)

    async function handleSubmit(e: React.FormEvent) {
        e.preventDefault()
        setError('')
        setLoading(true)
        try {
            const res = await apiFetch('/auth/verify-pin', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    pin,
                    deviceId: localStorage.getItem('deviceId'),
                }),
            })
            if (res.ok) {
                onSuccess()
            } else {
                setError('Incorrect PIN')
            }
        } catch {
            setError('Failed to reach server')
        } finally {
            setLoading(false)
        }
    }

    return (
        <div className="fixed inset-0 bg-black/70 flex items-center justify-center z-50">
            <form
                onSubmit={handleSubmit}
                className="bg-[#111] border border-[#1f1f1f] rounded-lg p-6 flex flex-col gap-3 w-72"
            >
                <h2 className="text-white text-lg">Enter PIN</h2>
                {error && <p className="text-[#e03030] text-sm">{error}</p>}
                <input
                    type="password"
                    inputMode="numeric"
                    autoFocus
                    value={pin}
                    onChange={e => setPin(e.target.value)}
                    className="bg-[#0a0a0a] border border-[#1f1f1f] text-white px-3 py-2 rounded-lg"
                    required
                />
                <div className="flex gap-2 justify-end">
                    <button
                        type="button"
                        onClick={onCancel}
                        className="px-3 py-2 text-sm text-[#888] hover:text-white"
                    >
                        Cancel
                    </button>
                    <button
                        type="submit"
                        disabled={loading}
                        className="bg-[#e03030] hover:bg-[#ff3c3c] text-white px-4 py-2 rounded-lg text-sm disabled:opacity-50"
                    >
                        {loading ? 'Checking...' : 'Confirm'}
                    </button>
                </div>
            </form>
        </div>
    )
}