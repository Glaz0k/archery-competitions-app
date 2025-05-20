// merited_master
// master_international
// master
// candidate_for_master
// first_class
// second_class
// third_class
import 'package:json_annotation/json_annotation.dart';

part 'common.g.dart';

@JsonEnum(fieldRename: FieldRename.snake)
enum SportsRank {
  meritedMaster,
  masterInternational,
  master,
  candidateForMaster,
  firstClass,
  secondClass,
  thirdClass;

  static List<String> stringValues = [
    "Заслуженный мастер спорта",
    "Мастер спорта международного класса",
    "Мастер спорта",
    "Кандидат в мастера спорта",
    "Первый разряд",
    "Второй разряд",
    "Третий разряд",
  ];

  factory SportsRank.fromString(String str) =>
      values[stringValues.indexOf(str)];

  @override
  String toString() => stringValues[index];
}

// male
// female
@JsonEnum()
enum Gender {
  male,
  female;

  @override
  String toString() {
    return switch (this) {
      Gender.male => "Мужчина",
      Gender.female => "Женщина",
    };
  }
}

// classic
// block
// classic_newbie
// 3D_classic
// 3D_compound
// 3D_long
// peripheral
// peripheral_with_ring
@JsonEnum()
enum BowClass {
  classic,
  block,
  @JsonValue("classic_newbie")
  classicNewbie,
  @JsonValue("3D_classic")
  classic3D,
  @JsonValue("3D_compound")
  compound3D,
  @JsonValue("3D_long")
  long3D,
  peripheral,
  @JsonValue("peripheral_with_ring")
  peripheralWithRing;

  static List<String> stringValues = [
    "Классический",
    "Блочный",
    "КЛ(новички)",
    "3Д-классический лук",
    "3Д-составной лук",
    "3Д-длинный лук",
    "Периферийный лук",
    "Периферийный лук(с кольцом)",
  ];

  factory BowClass.fromString(String str) => values[stringValues.indexOf(str)];

  @override
  String toString() => stringValues[index];
}

// created
// qualification_start
// qualification_end
// quarterfinal_start
// semifinal_start
// final_start
// completed
enum GroupState {
  created,
  qualificationStart,
  qualificationEnd,
  quarterfinalStart,
  semifinalStart,
  finalStart,
  completed,
}

// I
// II
// III
// F
enum CompetitionStage { I, II, III, F }

// ongoing
// top_win
// bot_win
@JsonEnum(fieldRename: FieldRename.snake)
enum SparringState { ongoing, topWin, botWin }

// "1-10"
// "6-10"
@JsonEnum()
enum RangeType {
  @JsonValue("1-10")
  one2ten,
  @JsonValue("6-10")
  six2ten,
}

// "shot_ordinal": <number>,
// "score": <string | null>
@JsonSerializable()
class Shot {
  int shotOrdinal;
  String? score;

  Shot(this.shotOrdinal, this.score);

  factory Shot.fromJson(Map<String, dynamic> json) => _$ShotFromJson(json);

  void changeBySlider(double sliderValue, RangeType type) {
    if (type == RangeType.six2ten && sliderValue < 6.0 || sliderValue == 0.0) {
      score = "M";
    } else {
      var temp = sliderValue.round();
      if (temp == 11) {
        score = "X";
      } else {
        score = temp.toString();
      }
    }
  }

  double toSliderValue(RangeType type) {
    return switch (score) {
      "M" => type == RangeType.six2ten ? 5.0 : 0.0,
      "X" => 11.0,
      null => type == RangeType.six2ten ? 8.0 : 5.0,
      _ => double.parse(score!),
    };
  }
}
