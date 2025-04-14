import createEnum from "../helper/createEnum";

const GroupState = createEnum({
  CREATED: { value: "created", textValue: "Создана" },
  QUAL_START: { value: "qualification_start", textValue: "Квалификация началась" },
  QUAL_END: { value: "qualification_end", textValue: "Квалификация закончилась" },
  QUARTERFINAL_START: { value: "quarterfinal_start", textValue: "Четвертьфинал начался" },
  SEMIFINAL_START: { value: "semifinal_start", textValue: "Полуфинал начался" },
  FINAL_START: { value: "final_start", textValue: "Финал начался" },
  COMPLETED: { value: "completed", textValue: "Закончена" },
});

export default GroupState;
