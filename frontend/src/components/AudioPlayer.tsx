// Audioplayer component to play audio files

import type { FileItem } from '../types'
import { API } from '../config'
import { formatSize, formatDate } from '../utils'
import { X } from 'lucide-react'

type Props = {
    file: FileItem
    onClose: () => void
}

export default function AudioPlayer({ file, onClose }: Props) {
    return (
        // Fixed bottom bar
        <div className = 'fixed bottom-0 left-0 right-0 bg-[#111] border-t border-[#1f1f1f] px-6 py-3 z-50'>
            <div className = 'flex items-center gap-6'>

                {/* Actual audio element */}
                <div className = 'flex flex-col gap-0.5 min-w-[160px]'>
                    <span className="text-sm font-medium text-white truncate">{file.name}</span>
                    <span className="text-xs text-[#55]">{formatSize(file.size)} · {formatDate(file.modTime)}</span>
                </div>

                {/* Audio controls */}
                <audio
                    controls
                    autoPlay
                    className="flex-1 h-8"
                    src={`${API}/file/stream?path=${encodeURIComponent(file.path)}`}
                />

                {/* Close button */}
                <button
                    onClick={onClose}
                    className="text-[#555] hover:text-white transition-colors duration-150"
                >
                    <X size={18} />
                </button>

            </div>

        </div>
    )
}