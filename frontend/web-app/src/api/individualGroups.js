import apiMock from "./mocks";

export async function deleteIndividualGroup(id) {
  console.log(id);
  await apiMock();
  return true;
}
