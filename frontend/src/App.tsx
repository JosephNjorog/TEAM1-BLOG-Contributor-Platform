import { Route, Routes } from "react-router-dom";
import type { Role } from "@team1/shared";
import { AppShell } from "./components/AppShell";
import { ProtectedRoute } from "./components/ProtectedRoute";
import { LoginPage } from "./pages/LoginPage";
import { RegisterPage } from "./pages/RegisterPage";
import { RoleHome } from "./pages/RoleHome";
import { ArticleEditorPage } from "./pages/contributor/ArticleEditorPage";

const ALL_ROLES: Role[] = ["super_admin", "moderator", "graphic_designer", "publisher", "contributor"];

export function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />

      <Route
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
      </Route>
    </Routes>
  );
}
