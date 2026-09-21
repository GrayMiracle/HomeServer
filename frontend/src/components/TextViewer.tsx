// Display text content
import { useState, useEffect } from 'react'
import type { FileItem } from '../types'
import { apiFetch } from '../api'
import Breadcrumb from './Breadcrumb'
import { X } from 'lucide-react'

type Props = {
    file: FileItem
    onClose: () => void
    onNavigate: (path: string) => void
}

export default function TextViewer({ file, onClose, onNavigate }: Props) {
    const [content, setContent] = useState('')
    const [error, setError] = useState('')
    const [loading, setLoading] = useState(false)

    useEffect(() => {
        async function fetchText() {
            setLoading(true)
            setError('')
            try {
                const res = await apiFetch(`/file/download?path=${encodeURIComponent(file.path)}&intent=view`)
                if (!res.ok) throw new Error('Failed to fetch')
                const text = await res.text()
                setContent(text)
                console.log('Fetched text content:', text)
            } catch {
                setError('Could not load file')
            } finally {
                setLoading(false)
            }
        }
        fetchText()
    }, [file.path])

    return (
        <div className="fixed inset-0 bg-[#0a0a0a] z-50 flex flex-col">
            {/* Top bar */}
            <div className="flex items-center justify-between px-6 py-3 bg-[#111] border-b border-[#1f1f1f] shrink-0">
                <Breadcrumb
                    currentPath={file.path}
                    onNavigate={(path) => {
                        onNavigate(path)
                        onClose()
                    }}
                />
                <button
                    onClick={onClose}
                    className="flex items-center gap-2 bg-[#1f1f1f] hover:bg-[#2a2a2a] text-white text-sm font-medium px-4 py-2 rounded-lg transition-colors duration-150"
                >
                    <X size={14} />
                </button>
            </div>
            {/* Content */}
            <div className="flex-1 overflow-auto p-6">
                {loading && <p className="text-sm text-[#666]">Loading...</p>}
                {error && <p className="text-sm text-[#e03030]">{error}</p>}
                {!loading && !error && (
                    <pre className="text-sm text-white font-mono whitespace-pre-wrap break-words leading-relaxed">
                        {content}
                    </pre>
                )}
            </div>
        </div>
    )
}