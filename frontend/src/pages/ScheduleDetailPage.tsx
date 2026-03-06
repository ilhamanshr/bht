import { useEffect, useState, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import type { Schedule } from '../types';
import { scheduleService, taskService } from '../services/api';
import { useGeolocation } from '../hooks/useGeolocation';
import { StatusBadge } from '../components/StatusBadge';
import { TaskItem } from '../components/TaskItem';

export function ScheduleDetailPage() {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const [schedule, setSchedule] = useState<Schedule | null>(null);
    const [loading, setLoading] = useState(true);
    const [actionLoading, setActionLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const { getLocation, loading: geoLoading } = useGeolocation();
    const [pendingAction, setPendingAction] = useState<'clock-in' | 'clock-out' | null>(null);
    const [showFallback, setShowFallback] = useState(false);

    const loadSchedule = useCallback(async () => {
        if (!id) return;
        try {
            setLoading(true);
            const data = await scheduleService.getById(parseInt(id));
            setSchedule(data);
        } catch (err) {
            console.error('Failed to load schedule:', err);
            setError('Schedule not found');
        } finally {
            setLoading(false);
        }
    }, [id]);

    useEffect(() => {
        loadSchedule();
    }, [loadSchedule]);

    const handleClockIn = async (skipGeo = false) => {
        if (!schedule) return;
        setActionLoading(true);
        setError(null);
        setShowFallback(false);
        try {
            let lat = null;
            let lng = null;

            if (!skipGeo) {
                try {
                    const location = await getLocation();
                    lat = location.latitude;
                    lng = location.longitude;
                } catch (err) {
                    console.warn('Geolocation failed, offering fallback:', err);
                    setPendingAction('clock-in');
                    setShowFallback(true);
                    setActionLoading(false);
                    return;
                }
            }

            const updated = await scheduleService.clockIn(schedule.id, lat as any, lng as any);
            setSchedule({ ...updated, tasks: schedule.tasks });
            await loadSchedule();
            setPendingAction(null);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to clock in');
            setPendingAction(null);
        } finally {
            setActionLoading(false);
        }
    };

    const handleClockOut = async (skipGeo = false) => {
        if (!schedule) return;
        setActionLoading(true);
        setError(null);
        setShowFallback(false);
        try {
            let lat = null;
            let lng = null;

            if (!skipGeo) {
                try {
                    const location = await getLocation();
                    lat = location.latitude;
                    lng = location.longitude;
                } catch (err) {
                    console.warn('Geolocation failed, offering fallback:', err);
                    setPendingAction('clock-out');
                    setShowFallback(true);
                    setActionLoading(false);
                    return;
                }
            }

            const updated = await scheduleService.clockOut(schedule.id, lat as any, lng as any);
            setSchedule({ ...updated, tasks: schedule.tasks });
            await loadSchedule();
            setPendingAction(null);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to clock out');
            setPendingAction(null);
        } finally {
            setActionLoading(false);
        }
    };

    const handleTaskUpdate = async (taskId: number, status: string, reason?: string) => {
        await taskService.update(taskId, status, reason);
        await loadSchedule();
    };

    const formatDateTime = (dt: string) => {
        return new Date(dt).toLocaleString('en-US', {
            month: 'short',
            day: 'numeric',
            hour: 'numeric',
            minute: '2-digit',
            hour12: true,
        });
    };

    if (loading) {
        return (
            <div className="page">
                <div className="loading">
                    <div className="loading__spinner"></div>
                    <p>Loading schedule details...</p>
                </div>
            </div>
        );
    }

    if (!schedule) {
        return (
            <div className="page">
                <div className="empty-state">
                    <span className="empty-state__icon">❌</span>
                    <p>Schedule not found</p>
                    <button className="btn btn--primary" onClick={() => navigate('/')}>
                        Back to Home
                    </button>
                </div>
            </div>
        );
    }

    const completedTasks = schedule.tasks?.filter((t) => t.status === 'completed').length || 0;
    const totalTasks = schedule.tasks?.length || 0;
    const progressPercentage = totalTasks > 0 ? Math.round((completedTasks / totalTasks) * 100) : 0;

    return (
        <div className="page">
            <header className="page__header">
                <div className="page__header-content">
                    <button className="btn btn--ghost btn--back" onClick={() => navigate('/')}>
                        ← Back
                    </button>
                    <h1 className="page__title">Schedule Details</h1>
                </div>
            </header>

            <main className="page__content">
                {error && (
                    <div className="alert alert--error">
                        <span>⚠️</span> {error}
                    </div>
                )}

                {showFallback && (
                    <div className="alert alert--warning alert--fallback">
                        <div className="alert__content">
                            <h4 className="alert__title">📍 Geolocation Failure</h4>
                            <p>Could not access your location. GPS is mandatory for verification.</p>
                            <div className="alert__actions">
                                <button
                                    className="btn btn--warning btn--sm"
                                    onClick={() => pendingAction === 'clock-in' ? handleClockIn(true) : handleClockOut(true)}
                                >
                                    Proceed & Flag for Review
                                </button>
                                <button className="btn btn--ghost btn--sm" onClick={() => { setShowFallback(false); setPendingAction(null); }}>
                                    Cancel
                                </button>
                            </div>
                        </div>
                    </div>
                )}

                {/* Schedule Info Card */}
                <div className="detail-card">
                    <div className="detail-card__header">
                        <h2 className="detail-card__client">{schedule.client_name}</h2>
                        <StatusBadge status={schedule.status} />
                    </div>

                    <div className="detail-card__info-grid">
                        <div className="detail-card__info">
                            <span className="detail-card__label">📅 Date</span>
                            <span className="detail-card__value">
                                {new Date(schedule.start_at).toLocaleDateString('en-US', {
                                    weekday: 'long',
                                    month: 'long',
                                    day: 'numeric',
                                    year: 'numeric',
                                })}
                            </span>
                        </div>
                        <div className="detail-card__info">
                            <span className="detail-card__label">🕐 Shift Time</span>
                            <span className="detail-card__value">
                                {new Date(schedule.start_at).toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit', hour12: true })} –{' '}
                                {new Date(schedule.end_at).toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit', hour12: true })}
                            </span>
                        </div>
                        <div className="detail-card__info">
                            <span className="detail-card__label">📍 Location</span>
                            <span className="detail-card__value">{schedule.location}</span>
                        </div>
                    </div>

                    {/* Clock In/Out Info */}
                    {schedule.clock_in_at && (
                        <div className="detail-card__clock">
                            <div className="detail-card__clock-item">
                                <span className="detail-card__clock-label">
                                    🟢 Clocked In
                                    {!schedule.clock_in_verified && (
                                        <span className="badge badge--warning ml-2" title="Unverified Location">⚠️ Unverified</span>
                                    )}
                                </span>
                                <span className="detail-card__clock-value">{formatDateTime(schedule.clock_in_at)}</span>
                                {schedule.clock_in_lat && schedule.clock_in_lng && (
                                    <span className="detail-card__clock-location">
                                        📍 {schedule.clock_in_lat.toFixed(4)}, {schedule.clock_in_lng.toFixed(4)}
                                    </span>
                                )}
                            </div>
                            {schedule.clock_out_at && (
                                <div className="detail-card__clock-item">
                                    <span className="detail-card__clock-label">
                                        🔴 Clocked Out
                                        {!schedule.clock_out_verified && (
                                            <span className="badge badge--warning ml-2" title="Unverified Location">⚠️ Unverified</span>
                                        )}
                                    </span>
                                    <span className="detail-card__clock-value">{formatDateTime(schedule.clock_out_at)}</span>
                                    {schedule.clock_out_lat && schedule.clock_out_lng && (
                                        <span className="detail-card__clock-location">
                                            📍 {schedule.clock_out_lat.toFixed(4)}, {schedule.clock_out_lng.toFixed(4)}
                                        </span>
                                    )}
                                </div>
                            )}
                        </div>
                    )}

                    {/* Action Buttons */}
                    <div className="detail-card__actions">
                        {schedule.status === 'upcoming' && (
                            <button
                                className="btn btn--primary btn--lg btn--full"
                                onClick={() => handleClockIn()}
                                disabled={actionLoading || geoLoading}
                            >
                                {actionLoading || geoLoading ? (
                                    <>
                                        <span className="btn__spinner"></span>
                                        {geoLoading ? 'Getting location...' : 'Starting visit...'}
                                    </>
                                ) : (
                                    '🟢 Start Visit'
                                )}
                            </button>
                        )}
                        {schedule.status === 'in_progress' && (
                            <button
                                className="btn btn--danger btn--lg btn--full"
                                onClick={() => handleClockOut()}
                                disabled={actionLoading || geoLoading}
                            >
                                {actionLoading || geoLoading ? (
                                    <>
                                        <span className="btn__spinner"></span>
                                        {geoLoading ? 'Getting location...' : 'Ending visit...'}
                                    </>
                                ) : (
                                    '🔴 End Visit'
                                )}
                            </button>
                        )}
                    </div>
                </div>

                {/* Task Progress */}
                {schedule.tasks && schedule.tasks.length > 0 && (
                    <div className="tasks-section">
                        <div className="tasks-section__header">
                            <h3 className="tasks-section__title">Care Activities</h3>
                            <span className="tasks-section__progress">
                                {completedTasks}/{totalTasks} completed
                            </span>
                        </div>

                        <div className="progress-bar">
                            <div
                                className="progress-bar__fill"
                                style={{ width: `${progressPercentage}%` }}
                            />
                        </div>

                        <div className="tasks-list">
                            {schedule.tasks.map((task) => (
                                <TaskItem
                                    key={task.id}
                                    task={task}
                                    canEdit={schedule.status === 'in_progress'}
                                    onUpdateStatus={handleTaskUpdate}
                                />
                            ))}
                        </div>
                    </div>
                )}
            </main>
        </div>
    );
}
