import { motion } from "framer-motion";

interface AnimatedListProps {
  children: React.ReactNode;
}

const listVariants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      staggerChildren: 0.1,
    },
  },
};

export default function AnimatedList({ children }: AnimatedListProps) {
  return (
    <motion.ul
      className="space-y-2"
      variants={listVariants}
      initial="hidden"
      animate="visible"
    >
      {children}
    </motion.ul>
  );
}
