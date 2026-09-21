// Imageviewer to display images

import type { FileItem } from '../types'
import { API } from '../config'
import { formatSize, formatDate } from '../utils'
import { X } from 'lucide-react'

// Props for img viewer
type Props = {
    file: FileItem
    onClose: () => void
}

// Default imgviewer function
export default function ImageViewer({ file, onClose }: Props) {
    return (
        // Fullscreen overlay, click outside to close
        <div
            className="fixed inset-0 bg-black/80 flex items-center justify-center z-50 p-6"
            onClick={onClose}
        >
            {/* Propogate click to prevent closing when in image container */}
            <div
                className="flex gap-6 w-full max-w-5xl items-start"
                onClick={e => e.stopPropagation()}
            >

                {/* Actual image */}
                <div className="flex-[3] flex items-center justify-center">
                    <img
                        className="max-w-full max-h-[80vh] rounded-lg object-contain"
                        src={`${API}/file/download?path=${encodeURIComponent(file.path)}`}
                        alt={file.name}
                    />
                </div>

                {/* Image extras */}
                <div className="flex-1 flex flex-col gap-4">
                    <button 
                        className="self-end text-[#555] hover:text-white transition-colors duration-150" 
                        onClick={onClose}
                    >
                        <X size={20} />
                    </button>
                    <h2 className="text-white font-medium text-base break-words">{file.name}</h2>
                    <div className="flex flex-col gap-2 text-sm">
                        <p className="text-[#55]">Size <span className="text-white ml-2">{formatSize(file.size)}</span></p>
                        <p className="text-[#55]">Modified <span className="text-white ml-2">{formatDate(file.modTime)}</span></p>
                    </div>
                </div>
            </div>
        </div>
    )
}

