// frontend/src/components/SettingsMenu.tsx
import { useState, useRef, useEffect } from 'react'
import { Settings, LogOut } from 'lucide-react'
import { logout } from '../api'

export default function SettingsMenu() {
    const [open, setOpen] = useState(false)
    const menuRef = useRef<HTMLDivElement>(null)

    // Close the menu if the user clicks anywhere outside it
    useEffect(() => {
        function handleClickOutside(e: MouseEvent) {
            if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
                setOpen(false)
            }
        }
        document.addEventListener('mousedown', handleClickOutside)
        return () => document.removeEventListener('mousedown', handleClickOutside)
    }, [])

    return (
        <div className="relative" ref={menuRef}>
            <button
                onClick={() => setOpen(prev => !prev)}
                className="p-2 rounded-lg hover:bg-[#161616] text-white transition-colors duration-150"
                aria-label="Settings"
            >
                <Settings size={20} />
            </button>

            {open && (
                <div className="absolute right-0 mt-2 w-44 bg-[#111] border border-[#1f1f1f] rounded-lg shadow-lg overflow-hidden z-50">
                    <button
                        onClick={logout}
                        className="w-full flex items-center gap-2 px-4 py-2 text-left text-white hover:bg-[#161616] transition-colors duration-150"
                    >
                        <LogOut size={16} className="text-[#e03030]" />
                        Log out
                    </button>
                </div>
            )}
        </div>
    )
}