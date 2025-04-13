export default async function apiMock(multi = 1) {
  const mockTime = 2000;
  await new Promise((res) => setTimeout(res, mockTime * multi));
}
