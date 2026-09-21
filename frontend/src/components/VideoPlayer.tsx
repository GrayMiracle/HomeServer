// VideoPlayer overlay for vid files
import type { FileItem } from '../types'
import { API } from '../config'
import { formatSize, formatDate } from '../utils'
import { X } from 'lucide-react'

type Props = {
    file: FileItem
    onClose: () => void
}

export default function VideoPlayer({ file, onClose }: Props) {
    return (
        <div
            className='fixed inset-0 bg-black/80 flex items-center justify-center z-50 p-6'
            onClick={onClose}
        >
            <div
                className='flex gap-6 w-full max-w-5xl max-h-[90vh] items-start overflow-y-auto bg-[#111] border border-[#1f1f1f] rounded-xl p-6'
                onClick={e => e.stopPropagation()}
            >
                {/* Actual video */}
                <div className="flex-[3] flex items-center justify-center">
                    <video
                        controls
                        autoPlay
                        className="max-h-[85vh] max-w-full w-auto rounded-lg"
                        src={`${API}/file/stream?path=${encodeURIComponent(file.path)}`}
                    />
                </div>

                {/* Video extras */}
                <div className="flex-1 flex flex-col gap-4">
                    <button
                        className="self-end text-[#555] hover:text-white transition-colors duration-150"
                        onClick={onClose}
                    >
                        <X size={20} />
                    </button>
                    <h2 className="text-white font-medium text-base break-words">{file.name}</h2>

                    <div className="flex flex-col gap-2 text-sm">
                        <p className="text-[#555]">Size <span className="text-white ml-2">{formatSize(file.size)}</span></p>
                        <p className="text-[#555]">Modified <span className="text-white ml-2">{formatDate(file.modTime)}</span></p>
                    </div>
                </div>
            </div>
        </div>
    )
}