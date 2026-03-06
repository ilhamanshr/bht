import { useEffect, useState } from 'react';
import type { Schedule, ScheduleStats } from '../types';
import { scheduleService } from '../services/api';
import { StatsCards } from '../components/StatsCards';
import { ScheduleCard } from '../components/ScheduleCard';

export function HomePage() {
    const [schedules, setSchedules] = useState<Schedule[]>([]);
    const [stats, setStats] = useState<ScheduleStats>({ total: 0, missed: 0, upcoming: 0, completed: 0 });
    const [loading, setLoading] = useState(true);
    const [filter, setFilter] = useState<'all' | 'today'>('all');

    useEffect(() => {
        loadData();
    }, [filter]);

    const loadData = async () => {
        try {
            setLoading(true);
            const [scheduleData, statsData] = await Promise.all([
                filter === 'today' ? scheduleService.getToday() : scheduleService.getAll(),
                scheduleService.getStats(),
            ]);
            setSchedules(scheduleData);
            setStats(statsData);
        } catch (err) {
            console.error('Failed to load data:', err);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="page">
            <header className="page__header">
                <div className="page__header-content">
                    <h1 className="page__title">
                        <span className="page__title-icon">🏥</span>
                        EVV Logger
                    </h1>
                    <p className="page__subtitle">Caregiver Shift Tracker</p>
                </div>
            </header>

            <main className="page__content">
                <StatsCards stats={stats} />

                <section className="schedule-section">
                    <div className="schedule-section__header">
                        <h2 className="schedule-section__title">Schedules</h2>
                        <div className="filter-tabs">
                            <button
                                className={`filter-tab ${filter === 'all' ? 'filter-tab--active' : ''}`}
                                onClick={() => setFilter('all')}
                            >
                                All
                            </button>
                            <button
                                className={`filter-tab ${filter === 'today' ? 'filter-tab--active' : ''}`}
                                onClick={() => setFilter('today')}
                            >
                                Today
                            </button>
                        </div>
                    </div>

                    {loading ? (
                        <div className="loading">
                            <div className="loading__spinner"></div>
                            <p>Loading schedules...</p>
                        </div>
                    ) : schedules.length === 0 ? (
                        <div className="empty-state">
                            <span className="empty-state__icon">📭</span>
                            <p>No schedules found</p>
                        </div>
                    ) : (
                        <div className="schedule-list">
                            {schedules.map((schedule) => (
                                <ScheduleCard key={schedule.id} schedule={schedule} />
                            ))}
                        </div>
                    )}
                </section>
            </main>
        </div>
    );
}
