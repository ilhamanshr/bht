import { useNavigate } from 'react-router-dom';
import type { Schedule } from '../types';
import { StatusBadge } from './StatusBadge';

interface ScheduleCardProps {
    schedule: Schedule;
}

export function ScheduleCard({ schedule }: ScheduleCardProps) {
    const navigate = useNavigate();

    return (
        <div
            className="schedule-card"
            onClick={() => navigate(`/schedules/${schedule.id}`)}
            role="button"
            tabIndex={0}
            onKeyDown={(e) => e.key === 'Enter' && navigate(`/schedules/${schedule.id}`)}
        >
            <div className="schedule-card__header">
                <h3 className="schedule-card__client">{schedule.client_name}</h3>
                <StatusBadge status={schedule.status} />
            </div>
            <div className="schedule-card__details">
                <div className="schedule-card__info">
                    <span className="schedule-card__icon">🕐</span>
                    <span>
                        {new Date(schedule.start_at).toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit', hour12: true })} –{' '}
                        {new Date(schedule.end_at).toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit', hour12: true })}
                    </span>
                </div>
                <div className="schedule-card__info">
                    <span className="schedule-card__icon">📍</span>
                    <span>{schedule.location}</span>
                </div>
                <div className="schedule-card__info">
                    <span className="schedule-card__icon">📅</span>
                    <span>{new Date(schedule.start_at).toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric' })}</span>
                </div>
            </div>
        </div>
    );
}
