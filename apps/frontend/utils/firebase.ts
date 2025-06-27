/**
 * Parses a Firebase error object and returns a user-friendly error message.
 * @param error - The Firebase error object.
 * @returns A string containing the parsed error message.
 */
export const parseFirebaseError = (error: any): string => {
  if (error?.code) {
    switch (error.code) {
      case "auth/invalid-email":
        return "Invalid email address format.";
      case "auth/user-disabled":
        return "This user account has been disabled.";
      case "auth/user-not-found":
        return "User not found. Please check your credentials.";
      case "auth/wrong-password":
        return "Invalid password. Please try again.";
      case "auth/email-already-in-use":
        return "This email address is already in use.";
      case "auth/operation-not-allowed":
        return "Email/password accounts are not enabled.";
      case "auth/weak-password":
        return "The password is too weak.";
      case "auth/invalid-credential":
        return "Invalid credentials provided.";
      default:
        return error.message || "An unknown error occurred.";
    }
  }
  return error.message || "An unknown error occurred.";
};
