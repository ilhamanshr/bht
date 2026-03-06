import { useState, useCallback } from 'react';
import type { GeoLocation } from '../types';

interface GeolocationState {
    location: GeoLocation | null;
    loading: boolean;
    error: string | null;
}

export function useGeolocation() {
    const [state, setState] = useState<GeolocationState>({
        location: null,
        loading: false,
        error: null,
    });

    const getLocation = useCallback((): Promise<GeoLocation> => {
        return new Promise((resolve, reject) => {
            setState({ location: null, loading: true, error: null });

            if (!navigator.geolocation) {
                const error = 'Geolocation is not supported by your browser';
                setState({ location: null, loading: false, error });
                reject(new Error(error));
                return;
            }

            navigator.geolocation.getCurrentPosition(
                (position) => {
                    const location: GeoLocation = {
                        latitude: position.coords.latitude,
                        longitude: position.coords.longitude,
                    };
                    setState({ location, loading: false, error: null });
                    resolve(location);
                },
                (err) => {
                    let errorMessage = 'Unable to retrieve your location';
                    switch (err.code) {
                        case err.PERMISSION_DENIED:
                            errorMessage = 'Location permission denied. Please enable location access in your browser settings.';
                            break;
                        case err.POSITION_UNAVAILABLE:
                            errorMessage = 'Location information unavailable. Using fallback.';
                            break;
                        case err.TIMEOUT:
                            errorMessage = 'Location request timed out. Please try again.';
                            break;
                    }
                    setState({ location: null, loading: false, error: errorMessage });
                    reject(new Error(errorMessage));
                },
                {
                    enableHighAccuracy: true,
                    timeout: 10000,
                    maximumAge: 0,
                }
            );
        });
    }, []);

    return { ...state, getLocation };
}
