import axios from 'axios';
import type { Schedule, ScheduleStats, Task, ApiResponse } from '../types';

const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL || '/api',
    headers: {
        'Content-Type': 'application/json',
    },
});

// Add timezone header to every request
api.interceptors.request.use((config) => {
    config.headers['X-Timezone'] = Intl.DateTimeFormat().resolvedOptions().timeZone;
    return config;
});

export const scheduleService = {
    getAll: async (): Promise<Schedule[]> => {
        const { data } = await api.get<ApiResponse<Schedule[]>>('/schedules');
        return data.data;
    },

    getToday: async (): Promise<Schedule[]> => {
        const { data } = await api.get<ApiResponse<Schedule[]>>('/schedules/today');
        return data.data;
    },

    getById: async (id: number): Promise<Schedule> => {
        const { data } = await api.get<ApiResponse<Schedule>>(`/schedules/${id}`);
        return data.data;
    },

    getStats: async (): Promise<ScheduleStats> => {
        const { data } = await api.get<ApiResponse<ScheduleStats>>('/schedules/stats');
        return data.data;
    },

    clockIn: async (id: number, latitude: number, longitude: number): Promise<Schedule> => {
        const { data } = await api.post<ApiResponse<Schedule>>(`/schedules/${id}/clock-in`, {
            latitude,
            longitude,
        });
        return data.data;
    },

    clockOut: async (id: number, latitude: number, longitude: number): Promise<Schedule> => {
        const { data } = await api.post<ApiResponse<Schedule>>(`/schedules/${id}/clock-out`, {
            latitude,
            longitude,
        });
        return data.data;
    },
};

export const taskService = {
    update: async (taskId: number, status: string, reason?: string): Promise<Task> => {
        const { data } = await api.post<ApiResponse<Task>>(`/tasks/${taskId}/update`, {
            status,
            reason: reason || null,
        });
        return data.data;
    },
};
