export default function (window) {
  const { protocol, hostname, port } = window.location;
  const host = port ? `${hostname}:${port}` : hostname;
  return { scheme: protocol.slice(0, -1), host };
}
