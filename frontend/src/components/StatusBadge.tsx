import type { Schedule } from '../types';

interface StatusBadgeProps {
    status: Schedule['status'] | string;
    size?: 'sm' | 'md';
}

export function StatusBadge({ status, size = 'md' }: StatusBadgeProps) {
    const getStatusConfig = (s: string) => {
        switch (s) {
            case 'completed':
                return { label: 'Completed', className: 'status-badge--completed' };
            case 'in_progress':
                return { label: 'In Progress', className: 'status-badge--in-progress' };
            case 'missed':
                return { label: 'Missed', className: 'status-badge--missed' };
            case 'not_completed':
                return { label: 'Not Completed', className: 'status-badge--missed' };
            case 'pending':
                return { label: 'Pending', className: 'status-badge--pending' };
            case 'upcoming':
            default:
                return { label: 'Upcoming', className: 'status-badge--upcoming' };
        }
    };

    const config = getStatusConfig(status);

    return (
        <span className={`status-badge ${config.className} ${size === 'sm' ? 'status-badge--sm' : ''}`}>
            {config.label}
        </span>
    );
}
