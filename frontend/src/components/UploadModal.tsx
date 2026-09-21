// Simple upload component
import { useState } from 'react'
import { apiFetch } from '../api'
import { X, Upload } from 'lucide-react'

// Props for the upload modal
type Props = {
    currentPath: string
    onUploadSuccess: () => void
    onClose: () => void
}

// Upload component function
export default function UploadModal({ currentPath, onUploadSuccess, onClose }: Props) {
    const [file, setFile] = useState<File | null>(null)
    const [error, setError] = useState('')
    const [uploading, setUploading] = useState(false)

    // function that handles uploads
    async function handleUpload() {
        if (!file) return
        setError('')
        setUploading(true)
        try {
            // formdata to add file as part of POST
            const formData = new FormData()
            formData.append('file', file)
            const res = await apiFetch(`/file/upload?path=${encodeURIComponent(currentPath)}`, {
                method: 'POST',
                body: formData
            })
            if (!res.ok) throw new Error('Upload failed')
            onUploadSuccess()
            onClose()
        } catch {
            setError('Upload failed')
        } finally {
            setUploading(false)
        }
    }

    // Main return
    return (
        // Backdrop
        <div 
            className="fixed inset-0 bg-black/80 flex items-center justify-center z-50"
            onClick = {onClose}
        >
            {/* Modal Box */}
            <div 
                className="bg-[#111] border border-[#1f1f1f] rounded-xl p-6 w-full max-w-md flex flex-col gap-4"
                onClick={e => e.stopPropagation()}
            >
                {/* Header */}
                <div className="flex items-center justify-between">
                    <h2 className="text-white font-medium text-base">Upload File</h2>
                    <button
                        onClick={onClose}
                        className="text-[#555] hover:text-white transition-colors duration-150"
                    >
                        <X size={18} />
                    </button>
                </div>

                {/* Destination */}
                <p className="text-xs text-[#555]">
                    Uploading to: <span className="text-white">{currentPath || 'Home'}</span>
                </p>

                {/* File input */}
                <input
                    type="file"
                    onChange={e => setFile(e.target.files ? e.target.files[0] : null)}
                    className="text-sm text-[#555] file:mr-4 file:py-2 file:px-4 file:rounded-lg file:border-0 file:bg-[#1f1f1f] file:text-white hover:file:bg-[#2a2a2a] file:cursor-pointer file:transition-colors file:duration-150"
                />

                {/* Error */}
                {error && <p className="text-xs text-[#e03030]">{error}</p>}

                {/* Upload button */}
                <button
                    onClick={handleUpload}
                    disabled={!file || uploading}
                    className="flex items-center justify-center gap-2 bg-[#e03030] hover:bg-[#ff3c3c] disabled:bg-[#2a2a2a] disabled:text-[#555] text-white text-sm font-medium px-4 py-2 rounded-lg transition-colors duration-150"
                >
                    <Upload size={14} />
                    {uploading ? 'Uploading...' : 'Upload'}
                </button>

            </div>
        </div>
    )

}