import type { ScheduleStats } from '../types';

interface StatsCardProps {
    stats: ScheduleStats;
}

export function StatsCards({ stats }: StatsCardProps) {
    const cards = [
        {
            label: 'Total Schedules',
            value: stats.total,
            icon: '📋',
            className: 'stats-card--total',
        },
        {
            label: 'Missed',
            value: stats.missed,
            icon: '❌',
            className: 'stats-card--missed',
        },
        {
            label: 'Upcoming Today',
            value: stats.upcoming,
            icon: '🕐',
            className: 'stats-card--upcoming',
        },
        {
            label: 'Completed Today',
            value: stats.completed,
            icon: '✅',
            className: 'stats-card--completed',
        },
    ];

    return (
        <div className="stats-grid">
            {cards.map((card) => (
                <div key={card.label} className={`stats-card ${card.className}`}>
                    <div className="stats-card__icon">{card.icon}</div>
                    <div className="stats-card__content">
                        <span className="stats-card__value">{card.value}</span>
                        <span className="stats-card__label">{card.label}</span>
                    </div>
                </div>
            ))}
        </div>
    );
}
