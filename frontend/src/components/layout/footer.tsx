export function Footer() {
  return (
    <footer className="border-t border-zinc-200 dark:border-zinc-800">
      <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-4 text-sm text-zinc-500 sm:px-6">
        <p>&copy; {new Date().getFullYear()} Creatrid</p>
        <p>Verified Creator Identity</p>
      </div>
    </footer>
  );
}
