import createEnum from "../../helper/createEnum";

const Gender = createEnum({
  MALE: { value: "male", textValue: "Мужчина" },
  FEMALE: { value: "female", textValue: "Женщина" },
});

export default Gender;
