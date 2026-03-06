import { useState } from 'react';
import type { Task } from '../types';
import { StatusBadge } from './StatusBadge';

interface TaskItemProps {
    task: Task;
    canEdit: boolean;
    onUpdateStatus: (taskId: number, status: string, reason?: string) => Promise<void>;
}

export function TaskItem({ task, canEdit, onUpdateStatus }: TaskItemProps) {
    const [showReasonInput, setShowReasonInput] = useState(false);
    const [reason, setReason] = useState('');
    const [loading, setLoading] = useState(false);

    const handleComplete = async () => {
        setLoading(true);
        try {
            await onUpdateStatus(task.id, 'completed');
        } finally {
            setLoading(false);
        }
    };

    const handleNotCompleted = () => {
        setShowReasonInput(true);
    };

    const handleSubmitReason = async () => {
        if (!reason.trim()) return;
        setLoading(true);
        try {
            await onUpdateStatus(task.id, 'not_completed', reason.trim());
            setShowReasonInput(false);
            setReason('');
        } finally {
            setLoading(false);
        }
    };

    const isCompleted = task.status === 'completed';
    const isNotCompleted = task.status === 'not_completed';
    const isPending = task.status === 'pending';

    return (
        <div className={`task-item ${isCompleted ? 'task-item--completed' : ''} ${isNotCompleted ? 'task-item--not-completed' : ''}`}>
            <div className="task-item__header">
                <div className="task-item__title-row">
                    <span className={`task-item__checkbox ${isCompleted ? 'task-item__checkbox--checked' : ''} ${isNotCompleted ? 'task-item__checkbox--failed' : ''}`}>
                        {isCompleted ? '✓' : isNotCompleted ? '✗' : ''}
                    </span>
                    <span className={`task-item__title ${isCompleted ? 'task-item__title--done' : ''}`}>
                        {task.title}
                    </span>
                </div>
                <StatusBadge status={task.status} size="sm" />
            </div>

            {isNotCompleted && task.reason && (
                <div className="task-item__reason">
                    <span className="task-item__reason-label">Reason:</span> {task.reason}
                </div>
            )}

            {canEdit && isPending && (
                <div className="task-item__actions">
                    {!showReasonInput ? (
                        <>
                            <button
                                className="btn btn--success btn--sm"
                                onClick={handleComplete}
                                disabled={loading}
                            >
                                {loading ? 'Saving...' : '✓ Complete'}
                            </button>
                            <button
                                className="btn btn--danger btn--sm"
                                onClick={handleNotCompleted}
                                disabled={loading}
                            >
                                ✗ Not Completed
                            </button>
                        </>
                    ) : (
                        <div className="task-item__reason-form">
                            <input
                                type="text"
                                className="task-item__reason-input"
                                placeholder="Enter reason..."
                                value={reason}
                                onChange={(e) => setReason(e.target.value)}
                                onKeyDown={(e) => e.key === 'Enter' && handleSubmitReason()}
                                autoFocus
                            />
                            <div className="task-item__reason-buttons">
                                <button
                                    className="btn btn--danger btn--sm"
                                    onClick={handleSubmitReason}
                                    disabled={loading || !reason.trim()}
                                >
                                    {loading ? 'Saving...' : 'Submit'}
                                </button>
                                <button
                                    className="btn btn--ghost btn--sm"
                                    onClick={() => { setShowReasonInput(false); setReason(''); }}
                                >
                                    Cancel
                                </button>
                            </div>
                        </div>
                    )}
                </div>
            )}
        </div>
    );
}
