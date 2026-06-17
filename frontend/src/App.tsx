import { Route, Routes } from "react-router-dom";
import type { Role } from "@team1/shared";
import { AppShell } from "./components/AppShell";
import { ProtectedRoute } from "./components/ProtectedRoute";
import { LandingPage } from "./pages/LandingPage";
import { LoginPage } from "./pages/LoginPage";
import { RegisterPage } from "./pages/RegisterPage";
import { RoleHome } from "./pages/RoleHome";
import { ArticleEditorPage } from "./pages/contributor/ArticleEditorPage";
import { PublishedHistoryPage } from "./pages/contributor/PublishedHistoryPage";
import { PaymentHistoryPage } from "./pages/contributor/PaymentHistoryPage";
import { ReviewArticlePage } from "./pages/moderator/ReviewArticlePage";
import { BannerUploadPage } from "./pages/designer/BannerUploadPage";
import { PublishArticlePage } from "./pages/publisher/PublishArticlePage";
import { PipelinePage } from "./pages/admin/PipelinePage";
import { ContributorsPage } from "./pages/admin/ContributorsPage";
import { PaymentsPage } from "./pages/admin/PaymentsPage";
import { UsersPage } from "./pages/admin/UsersPage";
import { AnalyticsPage } from "./pages/admin/AnalyticsPage";

const ALL_ROLES: Role[] = ["super_admin", "moderator", "graphic_designer", "publisher", "contributor"];

export function App() {
  return (
    <Routes>
      <Route path="/" element={<LandingPage />} />
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />

      <Route
        path="/app"
        element={
          <ProtectedRoute allow={ALL_ROLES}>
            <AppShell />
          </ProtectedRoute>
        }
      >
        <Route index element={<RoleHome />} />
        <Route
          path="articles/:id"
          element={
            <ProtectedRoute allow={["contributor"]}>
              <ArticleEditorPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="published"
          element={
            <ProtectedRoute allow={["contributor"]}>
              <PublishedHistoryPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="my-payments"
          element={
            <ProtectedRoute allow={["contributor"]}>
              <PaymentHistoryPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="review/:id"
          element={
            <ProtectedRoute allow={["moderator"]}>
              <ReviewArticlePage />
            </ProtectedRoute>
          }
        />
        <Route
          path="banner/:id"
          element={
            <ProtectedRoute allow={["graphic_designer"]}>
              <BannerUploadPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="publish/:id"
          element={
            <ProtectedRoute allow={["publisher"]}>
              <PublishArticlePage />
            </ProtectedRoute>
          }
        />
        <Route
          path="pipeline"
          element={
            <ProtectedRoute allow={["super_admin"]}>
              <PipelinePage />
            </ProtectedRoute>
          }
        />
        <Route
          path="contributors"
          element={
            <ProtectedRoute allow={["super_admin"]}>
              <ContributorsPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="payments"
          element={
            <ProtectedRoute allow={["super_admin"]}>
              <PaymentsPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="users"
          element={
            <ProtectedRoute allow={["super_admin"]}>
              <UsersPage />
            </ProtectedRoute>
          }
        />
        <Route
          path="analytics"
          element={
            <ProtectedRoute allow={["super_admin"]}>
              <AnalyticsPage />
            </ProtectedRoute>
          }
        />
      </Route>
    </Routes>
  );
}
