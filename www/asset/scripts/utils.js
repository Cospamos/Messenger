const isValidImageUrl = (url) => {
  if (typeof url !== 'string') return false;
  if (!url || url.trim() === '') return false;

  // Безопасная проверка расширения файла
  const safeExtensions = ['.jpg', '.jpeg', '.png', '.gif'];
  const safeProtocols = ['http://', 'https://'];
  const lowerUrl = url.toLowerCase();

  const sageExtensionTest = safeExtensions.some(ext => lowerUrl.endsWith(ext));
  const safeProtocolsTest = safeProtocols.some(ext => lowerUrl.startsWith(ext));

  return sageExtensionTest && safeProtocolsTest
}

