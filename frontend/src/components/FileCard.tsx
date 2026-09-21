// import global type structure
import { useRef } from 'react'
import type { FileItem } from '../types'
import { formatSize } from '../utils'
import { Trash2 } from 'lucide-react'

// Common filehandling type
type FileHandler = (file: FileItem) => void

// All filecard props
type Props = {
    file: FileItem
    onFolderClick: FileHandler
    onVideoClick: FileHandler
    onAudioClick: FileHandler
    onImageClick: FileHandler
    onPdfClick: FileHandler
    onTextClick: FileHandler
    onDownloadClick: FileHandler
    onDelete: FileHandler
}

const hoverSound = new Audio('/hover.mp3')

// Filecard component, clicks handled by file type within parent via props
export default function FileCard({ file, onFolderClick, onVideoClick, onAudioClick, onImageClick, onPdfClick, onTextClick, onDownloadClick, onDelete }: Props) {
    const cardRef = useRef<HTMLDivElement>(null)
    
    // Handle click based on file type
    function handleClick() {
        if (file.isDir) {
            onFolderClick(file)
        } else if (file.mimeType === "video") {
            onVideoClick(file)
        } else if (file.mimeType === "audio") {
            onAudioClick(file)
        } else if (file.mimeType === "image") {
            onImageClick(file)
        } else if (file.mimeType === 'pdf') {
            onPdfClick(file)
        } else if (file.mimeType === 'text') {
            onTextClick(file)
        } else {
            onDownloadClick(file)
        }
    }

    // Reset tilt on mouse leave
    function handleMouseLeave() {
        const card = cardRef.current
        if (!card) return
        card.style.transform = 'rotate(0deg) scale(1)'
    }

    // Play sound on hover
    function handleMouseEnter() {
        const card = cardRef.current
        if (!card) return
        card.style.transform = 'rotate(2deg) scale(1.3)'
        hoverSound.currentTime = 0
        hoverSound.volume = 0.3
        hoverSound.play()
    }

    // main file details within card
    return (
        <div 
            ref = {cardRef}
            onClick = {handleClick} 
            onMouseLeave = {handleMouseLeave}
            onMouseEnter = {handleMouseEnter}
            style = {{ transition: 'transform 0.1s ease, background 0.12s ease'}}
            className={`
                relative bg-[#111] border border-[#1f1f1f] rounded-lg p-4
                cursor-pointer select-none group
                hover:bg-[#161616] hover:border-[#2a2a2a] hover:z-10
                ${file.isDir
                    ? 'border-l-[3px] border-l-[#f0a500]'
                    : file.mimeType === 'video'  ? 'border-l-[3px] border-l-[#e03030]'
                    : file.mimeType === 'audio'  ? 'border-l-[3px] border-l-[#a855f7]'
                    : file.mimeType === 'image'  ? 'border-l-[3px] border-l-[#22c55e]'
                    : file.mimeType === 'pdf'    ? 'border-l-[3px] border-l-[#3b82f6]'
                    : file.mimeType === 'text'   ? 'border-l-[3px] border-l-[#888]'
                    :                             'border-l-[3px] border-l-[#444]'
                }
            `}
        >
            {/* File Name */}
            <p className="text-sm font-medium text-white truncate">{file.name}</p>

            {/* Metadata */}
            <p className="text-xs text-[#555] mt-1">
                {file.isDir ? 'Folder' : `${file.mimeType} · ${formatSize(file.size)}`}
            </p>

            {/* Delete */}
                <button
                onClick={(e) => { e.stopPropagation(); onDelete(file) }}
                className="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity duration-150"
            >
                <Trash2
                    size={14}
                    className="text-[#444] hover:text-[#e03030] transition-colors duration-300"
                />
            </button>
        </div>
    )
}
