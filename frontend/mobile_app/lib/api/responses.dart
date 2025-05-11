import 'package:json_annotation/json_annotation.dart';
import 'common.dart';

part 'responses.g.dart';

// "id": <number>,
// "full_name": <string>,
// "birth_date": <date ISO 8601>,
// "identity": <gender>,
// "bow": <bow_class | null>,
// "rank": <sports_rank | null>,
// "region": <string | null>,
// "federation": <string | null>,
// "club": <string | null>
@JsonSerializable(
  createToJson: true,
  fieldRename: FieldRename.snake,
  explicitToJson: true,
)
class CompetitorFull {
  int id;
  String fullName;
  DateTime birthDate;
  Gender identity;
  BowClass? bow;
  SportsRank? rank;
  String? region;
  String? federation;
  String? club;

  CompetitorFull(
    this.id,
    this.fullName,
    this.birthDate,
    this.identity,
    this.bow,
    this.rank,
    this.region,
    this.federation,
    this.club,
  );

  factory CompetitorFull.fromJson(Map<String, dynamic> json) =>
      _$CompetitorFullFromJson(json);
  Map<String, dynamic> toJson() => _$CompetitorFullToJson(this);

  CompetitorShrinked shrink() {
    return CompetitorShrinked(id, fullName);
  }
}

// "competition_id": <number>,
// "competitor": <competitor_full>,
// "is_active": <bool>,
// "created_at": <YYYY-MM-DDThh:mm:ss±hh ISO 8601>
@JsonSerializable(
  createToJson: true,
  fieldRename: FieldRename.snake,
  explicitToJson: true,
)
class CompetitorCompetitionDetail {
  int competitionId;
  CompetitorFull competitor;
  bool isActive;
  DateTime createdAt;

  CompetitorCompetitionDetail(
    this.competitionId,
    this.competitor,
    this.isActive,
    this.createdAt,
  );

  factory CompetitorCompetitionDetail.fromJson(Map<String, dynamic> json) =>
      _$CompetitorCompetitionDetailFromJson(json);
  Map<String, dynamic> toJson() => _$CompetitorCompetitionDetailToJson(this);
}

// "id": <number>,
// "competition_id": <number>,
// "bow": <bow_class>,
// "identity": <gender | null>,
// "state": <group_state>
@JsonSerializable(
  createToJson: true,
  fieldRename: FieldRename.snake,
  explicitToJson: true,
)
class IndividualGroup {
  int id;
  int competitionId;
  BowClass bow;
  Gender? identity;
  GroupState state;

  IndividualGroup(
    this.id,
    this.competitionId,
    this.bow,
    this.identity,
    this.state,
  );

  factory IndividualGroup.fromJson(Map<String, dynamic> json) =>
      _$IndividualGroupFromJson(json);
  Map<String, dynamic> toJson() => _$IndividualGroupToJson(this);
}

// "id": <number>,
// "title": <string>,
// "address": <string | null>,
// "season": <string | null>

@JsonSerializable(
  createToJson: true,
  fieldRename: FieldRename.snake,
  explicitToJson: true,
)
class Cup {
  int id;

  String title;

  String? address;

  String? season;

  Cup(this.id, this.title, this.address, this.season);

  factory Cup.fromJson(Map<String, dynamic> json) => _$CupFromJson(json);

  Map<String, dynamic> toJson() => _$CupToJson(this);
}

// "id": <number>,
// "stage": <competition_stage>,
// "start_date": <YYYY-MM-DD ISO 8601 | null>,
// "end_date": <YYYY-MM-DD ISO 8601 | null>,
// "is_ended": <bool>

@JsonSerializable()
class Competition {
  int id;
  CompetitionStage stage;
  DateTime? startDate;
  DateTime? endDate;
  bool isEnded;

  Competition(this.id, this.stage, this.startDate, this.endDate, this.isEnded);

  factory Competition.fromJson(Map<String, dynamic> json) => _$CompetitionFromJson(json);
}

// "group_id": <number>,
// "competitor": <competitor_full>
class CompetitorGroupDetail {
  int groupId;
  CompetitorFull competitor;

  CompetitorGroupDetail(this.groupId, this.competitor);
}

// "group_id": <number>,
// "distance": <string>,
// "round_count": <number>,
// "sections": [ <section> ]
class QualificationTable {
  int groupId;
  String distance;
  int roundCount;
  List<Section> sections;

  QualificationTable(
    this.groupId,
    this.distance,
    this.roundCount,
    this.sections,
  );
}

// "id": <number>,
// "competitor": <competitor_shrinked>,
// "place": <number | null>,
// "rounds": [ <qualification_round_shrinked> ],
// "total": <number | null>,
// "10_s": <number | null>,
// "9_s": <number | null>,
// "rank_gained": <sports_rank | null>
class Section {
  int id;
  CompetitorShrinked competitor;
  int? place;
  List<QualificationRoundShrinked> rounds;
  int? total;
  int? tenS;
  int? nineS;
  SportsRank? rankGained;

  Section(
    this.id,
    this.competitor,
    this.place,
    this.rounds,
    this.total,
    this.tenS,
    this.nineS,
    this.rankGained,
  );
}

// "id": <number>,
// "full_name": <string>
@JsonSerializable()
class CompetitorShrinked {
  int id;
  String fullName;

  CompetitorShrinked(this.id, this.fullName);
  factory CompetitorShrinked.fromJson(Map<String, dynamic> json) =>
      _$CompetitorShrinkedFromJson(json);
}

// "round_ordinal": <number>,
// "is_active": <bool>,
// "total_score": <number | null>
class QualificationRoundShrinked {
  int roundOrdinal;
  bool isActive;
  int? totalScore;

  QualificationRoundShrinked(this.roundOrdinal, this.isActive, this.totalScore);
}

// "group_id": <number>,
// "quarterfinal": <quarterfinal>,
// "semifinal": <semifinal | null>,
// "final": <final | null>
@JsonSerializable()
class FinalGrid {
  int groupId;
  Quarterfinal quarterfinal;
  Semifinal? semifinal;
  @JsonKey(name: "final")
  Final? fina1;

  FinalGrid(this.groupId, this.quarterfinal, this.semifinal, this.fina1);
  factory FinalGrid.fromJson(Map<String, dynamic> json) =>
      _$FinalGridFromJson(json);
}

// "sparring_1": <sparring | null>,
// "sparring_2": <sparring | null>,
// "sparring_3": <sparring | null>,
// "sparring_4": <sparring | null>
@JsonSerializable()
class Quarterfinal {
  Sparring? sparring1;
  Sparring? sparring2;
  Sparring? sparring3;
  Sparring? sparring4;

  Quarterfinal(this.sparring1, this.sparring2, this.sparring3, this.sparring4);
  factory Quarterfinal.fromJson(Map<String, dynamic> json) =>
      _$QuarterfinalFromJson(json);
}

// "sparring_5": <sparring | null>,
// "sparring_6": <sparring | null>
@JsonSerializable()
class Semifinal {
  Sparring? sparring5;
  Sparring? sparring6;

  Semifinal(this.sparring5, this.sparring6);
  factory Semifinal.fromJson(Map<String, dynamic> json) =>
      _$SemifinalFromJson(json);
}

// "sparring_gold": <sparring | null>,
// "sparring_bronze": <sparring | null>
@JsonSerializable()
class Final {
  Sparring? sparringGold;
  Sparring? sparringBronze;

  Final(this.sparringGold, this.sparringBronze);
  factory Final.fromJson(Map<String, dynamic> json) => _$FinalFromJson(json);
}

// "id": <number>,
// "top_place": <sparring_place | null>,
// "bot_place": <sparring_place | null>,
// "state": <sparring_state>
@JsonSerializable()
class Sparring {
  int id;
  SparringPlace? topPlace;
  SparringPlace? botPlace;
  SparringState state;

  Sparring(this.id, this.topPlace, this.botPlace, this.state);
  factory Sparring.fromJson(Map<String, dynamic> json) =>
      _$SparringFromJson(json);
}

// "id": <number>,
// "competitor": <competitor_shrinked>,
// "range_group": <range_group>,
// "is_active": <bool>,
// "shoot_out": <shoot_out | null>,
// "sparring_score": <number>
@JsonSerializable()
class SparringPlace {
  int id;
  CompetitorShrinked competitor;
  RangeGroup rangeGroup;
  bool isActive;
  ShootOut? shootOut;
  int sparringScore;

  SparringPlace(
    this.id,
    this.competitor,
    this.rangeGroup,
    this.isActive,
    this.shootOut,
    this.sparringScore,
  );
  factory SparringPlace.fromJson(Map<String, dynamic> json) =>
      _$SparringPlaceFromJson(json);
}

// "id": <number>,
// "ranges_max_count": <number>,
// "range_size": <number>,
// "type": <range_type>
// "ranges": <[ <range> ]>,
// "total_score": <number | null>
@JsonSerializable()
class RangeGroup {
  int id;
  int rangesMaxCount;
  int rangeSize;
  RangeType type;
  List<Range> ranges;
  int? totalScore;

  RangeGroup(
    this.id,
    this.rangesMaxCount,
    this.rangeSize,
    this.type,
    this.ranges,
    this.totalScore,
  );
  factory RangeGroup.fromJson(Map<String, dynamic> json) =>
      _$RangeGroupFromJson(json);
}

// "id": <number>,
// "range_ordinal": <number>,
// "is_active": <bool>,
// "shots": <[ <shot> ] | null>,
// "range_score": <number | null>
@JsonSerializable()
class Range {
  int id;
  int rangeOrdinal;
  bool isActive;
  List<Shot>? shots;
  int? rangeScore;

  Range(this.id, this.rangeOrdinal, this.isActive, this.shots, this.rangeScore);
  factory Range.fromJson(Map<String, dynamic> json) => _$RangeFromJson(json);
}

// id": <number>,
// "score": <string | null>,
// "priority": <bool | null>
@JsonSerializable()
class ShootOut {
  int id;
  String? score;
  bool? priority;

  ShootOut(this.id, this.score, this.priority);
  factory ShootOut.fromJson(Map<String, dynamic> json) =>
      _$ShootOutFromJson(json);
}

// "section_id": <number>,
// "round_ordinal": <number>,
// "is_active": <bool>,
// "range_group": <range_group>
class QualificationRoundFull {
  int sectionId;
  int roundOrdinal;
  bool isActive;
  RangeGroup rangeGroup;

  QualificationRoundFull(
    this.sectionId,
    this.roundOrdinal,
    this.isActive,
    this.rangeGroup,
  );
}

// "id": <number>,
// "competitor": <competitor_shrinked>,
// "range_group": <range_group>,
// "is_active": <bool>,
// "shoot_out": <shoot_out | null>,
// "sparring_score": <number>
@JsonSerializable()
class SparingPlace {
  int id;
  CompetitorShrinked competitor;
  RangeGroup rangeGroup;
  bool isActive;
  ShootOut? shootOut;
  int sparringScore;

  SparingPlace(
    this.id,
    this.competitor,
    this.rangeGroup,
    this.isActive,
    this.shootOut,
    this.sparringScore,
  );
  factory SparingPlace.fromJson(Map<String, dynamic> json) =>
      _$SparingPlaceFromJson(json);
}
