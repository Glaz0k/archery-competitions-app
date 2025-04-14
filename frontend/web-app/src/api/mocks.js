export default async function apiMock(multi = 1) {
  const mockTime = 1000;
  await new Promise((res) => setTimeout(res, mockTime * multi));
}
