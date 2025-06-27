import { motion } from "framer-motion";

interface AnimatedListItemProps {
  children: React.ReactNode;
  index: number;
}

export default function AnimatedListItem({
  children,
  index,
}: AnimatedListItemProps) {
  return (
    <motion.li
      layout="position"
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3, delay: index * 0.05 }}
    >
      {children}
    </motion.li>
  );
}
