import { initializeApp, getApps, getApp } from "firebase/app";
import {
  GoogleAuthProvider,
  signInWithPopup,
  signInAnonymously,
  signOut as firebaseSignOut,
  onAuthStateChanged,
} from "firebase/auth";
import { getAuth } from "firebase/auth";

/**
 * Firebase configuration object.
 * @see https://firebase.google.com/docs/web/setup#config-object
 */
const firebaseConfig = {
  apiKey: process.env.NEXT_PUBLIC_FIREBASE_API_KEY,
  authDomain: process.env.NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN,
  projectId: process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID,
  storageBucket: process.env.NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET,
  messagingSenderId: process.env.NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID,
  appId: process.env.NEXT_PUBLIC_FIREBASE_APP_ID,
};

/**
 * Initializes Firebase, creating a new app instance if one doesn't already exist.
 * @see https://firebase.google.com/docs/web/setup#initialize-firebase
 */
const app = !getApps().length ? initializeApp(firebaseConfig) : getApp();

/**
 * Firebase Authentication instance.
 * @see https://firebase.google.com/docs/auth
 */
const auth = getAuth(app);

export { app, auth, onAuthStateChanged };

/**
 * Google Authentication provider.
 * @see https://firebase.google.com/docs/auth/web/google-signin
 */
const provider = new GoogleAuthProvider();

/**
 * Signs in the user with Google.
 * @returns A promise that resolves with the user's ID token.
 * @throws An error if the sign-in process fails.
 */
export const signInWithGoogle = async () => {
  try {
    const result = await signInWithPopup(auth, provider);
    const idToken = await result.user.getIdToken();
    return idToken;
  } catch (error) {
    console.error("Error during Google sign-in:", error);
    throw error;
  }
};

/**
 * Signs in the user anonymously.
 * @returns A promise that resolves with the user's ID token.
 * @throws An error if the sign-in process fails.
 */
export const signInAsGuest = async () => {
  try {
    const userCredential = await signInAnonymously(auth);
    const idToken = await userCredential.user.getIdToken();
    return idToken;
  } catch (error) {
    console.error("Error during guest sign-in:", error);
    throw error;
  }
};

/**
 * Signs out the current user.
 * @throws An error if the sign-out process fails.
 */
export const signOut = async () => {
  try {
    await firebaseSignOut(auth);
  } catch (error) {
    console.error("Error signing out:", error);
    throw error;
  }
};
