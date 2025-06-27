"use client";

import { parseFirebaseError } from "@/utils/firebase";
import {
  createUserWithEmailAndPassword,
  signInWithEmailAndPassword,
} from "firebase/auth";
import { useRouter } from "next/navigation";
import { useState, useTransition } from "react";
import { useAuth as useAuthContext } from "../context/AuthContext";
import {
  auth,
  signInAsGuest as firebaseSignInAsGuest,
  signInWithGoogle as firebaseSignInWithGoogle,
  signOut as firebaseSignOut,
} from "../lib/firebase";
import * as authService from "../services/api/authService";

/**
 * Custom hook for authentication-related actions.
 * Provides functions for signing in with Google, signing out, and accessing user context.
 * @returns An object containing authentication status, user data, and action functions.
 */
export const useAuth = () => {
  const context = useAuthContext();
  const router = useRouter();

  const [error, setError] = useState<string | null>(null);
  const [loading, startTransition] = useTransition();

  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }

  /**
   * Initiates the Google sign-in process.
   * On success, signs up the user in the backend and redirects to the dashboard.
   */
  const signInWithGoogle = async () => {
    startTransition(async () => {
      try {
        const idToken = await firebaseSignInWithGoogle();

        if (idToken) {
          await authService.signUp(idToken);
          router.push("/dashboard");
        }
      } catch (error: any) {
        setError(parseFirebaseError(error));
      }
    });
  };

  /**
   * Initiates the guest sign-in process.
   * On success, signs up the user in the backend and redirects to the dashboard.
   */
  const signInAsGuest = async () => {
    startTransition(async () => {
      try {
        const idToken = await firebaseSignInAsGuest();

        if (idToken) {
          await authService.signUp(idToken);
          router.push("/dashboard");
        }
      } catch (error: any) {
        setError(parseFirebaseError(error));
      }
    });
  };

  /**
   * Initiates the email/password sign-in process.
   * On success, signs up the user in the backend and redirects to the dashboard.
   * @param email - The user's email.
   * @param password - The user's password.
   * @param type - The type of action to perform, either "create" or "login".
   */
  const signInAsEmail = async (
    email: string,
    password: string,
    type: "create" | "login"
  ) => {
    startTransition(async () => {
      try {
        const userCredential =
          type === "create"
            ? await createUserWithEmailAndPassword(auth, email, password)
            : await signInWithEmailAndPassword(auth, email, password);

        const user = userCredential.user;
        const idToken = await user.getIdToken();

        if (idToken) {
          await authService.signUp(idToken);
          router.push("/dashboard");
        }
      } catch (error: any) {
        setError(parseFirebaseError(error));
      }
    });
  };

  /**
   * Signs the current user out.
   * On success, redirects to the login page.
   */
  const signOut = async () => {
    await firebaseSignOut();
    router.push("/login");
  };

  return {
    ...context,
    signInWithGoogle,
    signInAsGuest,
    signInAsEmail,
    signOut,
    error,
    signUp: authService.signUp,
    loading,
  };
};
