export interface Schedule {
  id: number;
  client_name: string;
  start_at: string;
  end_at: string;
  location: string;
  status: 'upcoming' | 'in_progress' | 'completed' | 'missed';
  clock_in_at: string | null;
  clock_in_lat: number | null;
  clock_in_lng: number | null;
  clock_out_at: string | null;
  clock_out_lat: number | null;
  clock_out_lng: number | null;
  tasks?: Task[];
  created_at: string;
  updated_at: string;
  clock_in_verified: boolean;
  clock_out_verified: boolean;
}

export interface Task {
  id: number;
  schedule_id: number;
  title: string;
  status: 'pending' | 'completed' | 'not_completed';
  reason: string | null;
  created_at: string;
  updated_at: string;
}

export interface ScheduleStats {
  total: number;
  missed: number;
  upcoming: number;
  completed: number;
}

export interface GeoLocation {
  latitude: number;
  longitude: number;
}

export interface ApiResponse<T> {
  data: T;
}

export interface ApiError {
  error: string;
}
