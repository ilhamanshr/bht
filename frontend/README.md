# Mini EVV Logger - Frontend

This is the frontend client for the Mini EVV Logger application, a comprehensive Caregiver Shift Tracker designed for Electronic Visit Verification (EVV) compliance.

## Tech Stack
- **Framework:** React 19
- **Language:** TypeScript
- **Build Tool:** Vite
- **Styling:** CSS Modules / Custom CSS Variables (Dark Theme)
- **Icons:** Lucide React
- **Date Handling:** `date-fns` & `date-fns-tz`

## Key Decisions
- **Vite over Create React App**: Chosen for its significantly faster HMR (Hot Module Replacement) and optimized build speeds.
- **Custom CSS Variables**: Instead of a heavy component library, we used raw CSS with custom properties (`var(--primary)`, `var(--bg-dark)`) to guarantee a lightweight footprint and complete control over the bespoke dark mode aesthetic.
- **Geolocation API**: We utilize the native browser `navigator.geolocation` API to ensure EVV compliance when clocking in and out of visits.
- **Timezone Awareness**: The frontend passes an `X-Timezone` header (e.g., `Asia/Tokyo`, `America/New_York`) with every request using an Axios interceptor, ensuring the backend calculates "today's" schedule relative to the caregiver's actual physical location.

## Local Setup Instructions

### Prerequisites
- Node.js 22+
- npm or yarn

### 1. Environment Variables
Create a `.env` file in the `frontend` directory (you can copy `.env.example` if it exists):
```env
VITE_API_BASE_URL=http://localhost:8080/api
```

### 2. Install Dependencies
```bash
npm install
```

### 3. Run Development Server
```bash
npm run dev
```
The application will be accessible at `http://localhost:5173`. 
*(Note: Vite proxy is configured to forward `/api` requests to `http://localhost:8080` to prevent CORS issues during local development).*

## Assumptions
- The application currently operates in a single-user mode for demonstration purposes (no JWT login screen).
- The browser supports the Geolocation API.

## Future Improvements
- Implement state management (Zustand or Redux) if the application scales further.
- Add comprehensive frontend tests using Vitest and React Testing Library.
- Implement a PWA (Progressive Web App) service worker to allow caregivers to view schedules offline.
