// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'responses.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

CompetitorFull _$CompetitorFullFromJson(Map<String, dynamic> json) =>
    CompetitorFull(
      (json['id'] as num).toInt(),
      json['full_name'] as String,
      DateTime.parse(json['birth_date'] as String),
      $enumDecode(_$GenderEnumMap, json['identity']),
      $enumDecodeNullable(_$BowClassEnumMap, json['bow']),
      $enumDecodeNullable(_$SportsRankEnumMap, json['rank']),
      json['region'] as String?,
      json['federation'] as String?,
      json['club'] as String?,
    );

Map<String, dynamic> _$CompetitorFullToJson(CompetitorFull instance) =>
    <String, dynamic>{
      'id': instance.id,
      'full_name': instance.fullName,
      'birth_date': instance.birthDate.toIso8601String(),
      'identity': _$GenderEnumMap[instance.identity]!,
      'bow': _$BowClassEnumMap[instance.bow],
      'rank': _$SportsRankEnumMap[instance.rank],
      'region': instance.region,
      'federation': instance.federation,
      'club': instance.club,
    };

const _$GenderEnumMap = {Gender.male: 'male', Gender.female: 'female'};

const _$BowClassEnumMap = {
  BowClass.classic: 'classic',
  BowClass.block: 'block',
  BowClass.classicNewbie: 'classic_newbie',
  BowClass.classic3D: '3D_classic',
  BowClass.compound3D: '3D_compound',
  BowClass.long3D: '3D_long',
  BowClass.peripheral: 'peripheral',
  BowClass.peripheralWithRing: 'peripheral_with_ring',
};

const _$SportsRankEnumMap = {
  SportsRank.meritedMaster: 'merited_master',
  SportsRank.masterInternational: 'master_international',
  SportsRank.master: 'master',
  SportsRank.candidateForMaster: 'candidate_for_master',
  SportsRank.firstClass: 'first_class',
  SportsRank.secondClass: 'second_class',
  SportsRank.thirdClass: 'third_class',
};

CompetitorCompetitionDetail _$CompetitorCompetitionDetailFromJson(
  Map<String, dynamic> json,
) => CompetitorCompetitionDetail(
  (json['competition_id'] as num).toInt(),
  CompetitorFull.fromJson(json['competitor'] as Map<String, dynamic>),
  json['is_active'] as bool,
  DateTime.parse(json['created_at'] as String),
);

Map<String, dynamic> _$CompetitorCompetitionDetailToJson(
  CompetitorCompetitionDetail instance,
) => <String, dynamic>{
  'competition_id': instance.competitionId,
  'competitor': instance.competitor.toJson(),
  'is_active': instance.isActive,
  'created_at': instance.createdAt.toIso8601String(),
};

IndividualGroup _$IndividualGroupFromJson(Map<String, dynamic> json) =>
    IndividualGroup(
      (json['id'] as num).toInt(),
      (json['competition_id'] as num).toInt(),
      $enumDecode(_$BowClassEnumMap, json['bow']),
      $enumDecodeNullable(_$GenderEnumMap, json['identity']),
      $enumDecode(_$GroupStateEnumMap, json['state']),
    );

Map<String, dynamic> _$IndividualGroupToJson(IndividualGroup instance) =>
    <String, dynamic>{
      'id': instance.id,
      'competition_id': instance.competitionId,
      'bow': _$BowClassEnumMap[instance.bow]!,
      'identity': _$GenderEnumMap[instance.identity],
      'state': _$GroupStateEnumMap[instance.state]!,
    };

const _$GroupStateEnumMap = {
  GroupState.created: 'created',
  GroupState.qualificationStart: 'qualificationStart',
  GroupState.qualificationEnd: 'qualificationEnd',
  GroupState.quarterfinalStart: 'quarterfinalStart',
  GroupState.semifinalStart: 'semifinalStart',
  GroupState.finalStart: 'finalStart',
  GroupState.completed: 'completed',
};

Cup _$CupFromJson(Map<String, dynamic> json) => Cup(
  (json['id'] as num).toInt(),
  json['title'] as String,
  json['address'] as String?,
  json['season'] as String?,
);

Map<String, dynamic> _$CupToJson(Cup instance) => <String, dynamic>{
  'id': instance.id,
  'title': instance.title,
  'address': instance.address,
  'season': instance.season,
};

Competition _$CompetitionFromJson(Map<String, dynamic> json) => Competition(
  (json['id'] as num).toInt(),
  $enumDecode(_$CompetitionStageEnumMap, json['stage']),
  json['start_date'] == null
      ? null
      : DateTime.parse(json['start_date'] as String),
  json['end_date'] == null ? null : DateTime.parse(json['end_date'] as String),
  json['is_ended'] as bool,
);

Map<String, dynamic> _$CompetitionToJson(Competition instance) =>
    <String, dynamic>{
      'id': instance.id,
      'stage': _$CompetitionStageEnumMap[instance.stage]!,
      'start_date': instance.startDate?.toIso8601String(),
      'end_date': instance.endDate?.toIso8601String(),
      'is_ended': instance.isEnded,
    };

const _$CompetitionStageEnumMap = {
  CompetitionStage.I: 'I',
  CompetitionStage.II: 'II',
  CompetitionStage.III: 'III',
  CompetitionStage.F: 'F',
};

CompetitorShrinked _$CompetitorShrinkedFromJson(Map<String, dynamic> json) =>
    CompetitorShrinked(
      (json['id'] as num).toInt(),
      json['full_name'] as String,
    );

Map<String, dynamic> _$CompetitorShrinkedToJson(CompetitorShrinked instance) =>
    <String, dynamic>{'id': instance.id, 'full_name': instance.fullName};

FinalGrid _$FinalGridFromJson(Map<String, dynamic> json) => FinalGrid(
  (json['group_id'] as num).toInt(),
  Quarterfinal.fromJson(json['quarterfinal'] as Map<String, dynamic>),
  json['semifinal'] == null
      ? null
      : Semifinal.fromJson(json['semifinal'] as Map<String, dynamic>),
  json['final'] == null
      ? null
      : Final.fromJson(json['final'] as Map<String, dynamic>),
);

Map<String, dynamic> _$FinalGridToJson(FinalGrid instance) => <String, dynamic>{
  'group_id': instance.groupId,
  'quarterfinal': instance.quarterfinal,
  'semifinal': instance.semifinal,
  'final': instance.fina1,
};

Quarterfinal _$QuarterfinalFromJson(Map<String, dynamic> json) => Quarterfinal(
  json['sparring1'] == null
      ? null
      : Sparring.fromJson(json['sparring1'] as Map<String, dynamic>),
  json['sparring2'] == null
      ? null
      : Sparring.fromJson(json['sparring2'] as Map<String, dynamic>),
  json['sparring3'] == null
      ? null
      : Sparring.fromJson(json['sparring3'] as Map<String, dynamic>),
  json['sparring4'] == null
      ? null
      : Sparring.fromJson(json['sparring4'] as Map<String, dynamic>),
);

Map<String, dynamic> _$QuarterfinalToJson(Quarterfinal instance) =>
    <String, dynamic>{
      'sparring1': instance.sparring1,
      'sparring2': instance.sparring2,
      'sparring3': instance.sparring3,
      'sparring4': instance.sparring4,
    };

Semifinal _$SemifinalFromJson(Map<String, dynamic> json) => Semifinal(
  json['sparring5'] == null
      ? null
      : Sparring.fromJson(json['sparring5'] as Map<String, dynamic>),
  json['sparring6'] == null
      ? null
      : Sparring.fromJson(json['sparring6'] as Map<String, dynamic>),
);

Map<String, dynamic> _$SemifinalToJson(Semifinal instance) => <String, dynamic>{
  'sparring5': instance.sparring5,
  'sparring6': instance.sparring6,
};

Final _$FinalFromJson(Map<String, dynamic> json) => Final(
  json['sparring_gold'] == null
      ? null
      : Sparring.fromJson(json['sparring_gold'] as Map<String, dynamic>),
  json['sparring_bronze'] == null
      ? null
      : Sparring.fromJson(json['sparring_bronze'] as Map<String, dynamic>),
);

Map<String, dynamic> _$FinalToJson(Final instance) => <String, dynamic>{
  'sparring_gold': instance.sparringGold,
  'sparring_bronze': instance.sparringBronze,
};

Sparring _$SparringFromJson(Map<String, dynamic> json) => Sparring(
  (json['id'] as num).toInt(),
  json['top_place'] == null
      ? null
      : SparringPlace.fromJson(json['top_place'] as Map<String, dynamic>),
  json['bot_place'] == null
      ? null
      : SparringPlace.fromJson(json['bot_place'] as Map<String, dynamic>),
  $enumDecode(_$SparringStateEnumMap, json['state']),
);

Map<String, dynamic> _$SparringToJson(Sparring instance) => <String, dynamic>{
  'id': instance.id,
  'top_place': instance.topPlace,
  'bot_place': instance.botPlace,
  'state': _$SparringStateEnumMap[instance.state]!,
};

const _$SparringStateEnumMap = {
  SparringState.ongoing: 'ongoing',
  SparringState.topWin: 'top_win',
  SparringState.botWin: 'bot_win',
};

SparringPlace _$SparringPlaceFromJson(Map<String, dynamic> json) =>
    SparringPlace(
      (json['id'] as num).toInt(),
      CompetitorShrinked.fromJson(json['competitor'] as Map<String, dynamic>),
      RangeGroup.fromJson(json['range_group'] as Map<String, dynamic>),
      json['is_active'] as bool,
      json['shoot_out'] == null
          ? null
          : ShootOut.fromJson(json['shoot_out'] as Map<String, dynamic>),
      (json['sparring_score'] as num).toInt(),
    );

Map<String, dynamic> _$SparringPlaceToJson(SparringPlace instance) =>
    <String, dynamic>{
      'id': instance.id,
      'competitor': instance.competitor,
      'range_group': instance.rangeGroup,
      'is_active': instance.isActive,
      'shoot_out': instance.shootOut,
      'sparring_score': instance.sparringScore,
    };

RangeGroup _$RangeGroupFromJson(Map<String, dynamic> json) => RangeGroup(
  (json['id'] as num).toInt(),
  (json['ranges_max_count'] as num).toInt(),
  (json['range_size'] as num).toInt(),
  $enumDecode(_$RangeTypeEnumMap, json['type']),
  (json['ranges'] as List<dynamic>)
      .map((e) => Range.fromJson(e as Map<String, dynamic>))
      .toList(),
  (json['total_score'] as num?)?.toInt(),
);

Map<String, dynamic> _$RangeGroupToJson(RangeGroup instance) =>
    <String, dynamic>{
      'id': instance.id,
      'ranges_max_count': instance.rangesMaxCount,
      'range_size': instance.rangeSize,
      'type': _$RangeTypeEnumMap[instance.type]!,
      'ranges': instance.ranges,
      'total_score': instance.totalScore,
    };

const _$RangeTypeEnumMap = {
  RangeType.one2ten: '1-10',
  RangeType.six2ten: '6-10',
};

Range _$RangeFromJson(Map<String, dynamic> json) => Range(
  (json['id'] as num).toInt(),
  (json['range_ordinal'] as num).toInt(),
  json['is_active'] as bool,
  (json['shots'] as List<dynamic>?)
      ?.map((e) => Shot.fromJson(e as Map<String, dynamic>))
      .toList(),
  (json['range_score'] as num?)?.toInt(),
);

Map<String, dynamic> _$RangeToJson(Range instance) => <String, dynamic>{
  'id': instance.id,
  'range_ordinal': instance.rangeOrdinal,
  'is_active': instance.isActive,
  'shots': instance.shots,
  'range_score': instance.rangeScore,
};

ShootOut _$ShootOutFromJson(Map<String, dynamic> json) => ShootOut(
  (json['id'] as num).toInt(),
  json['score'] as String?,
  json['priority'] as bool?,
);

Map<String, dynamic> _$ShootOutToJson(ShootOut instance) => <String, dynamic>{
  'id': instance.id,
  'score': instance.score,
  'priority': instance.priority,
};

SparingPlace _$SparingPlaceFromJson(Map<String, dynamic> json) => SparingPlace(
  (json['id'] as num).toInt(),
  CompetitorShrinked.fromJson(json['competitor'] as Map<String, dynamic>),
  RangeGroup.fromJson(json['range_group'] as Map<String, dynamic>),
  json['is_active'] as bool,
  json['shoot_out'] == null
      ? null
      : ShootOut.fromJson(json['shoot_out'] as Map<String, dynamic>),
  (json['sparring_score'] as num).toInt(),
);

Map<String, dynamic> _$SparingPlaceToJson(SparingPlace instance) =>
    <String, dynamic>{
      'id': instance.id,
      'competitor': instance.competitor,
      'range_group': instance.rangeGroup,
      'is_active': instance.isActive,
      'shoot_out': instance.shootOut,
      'sparring_score': instance.sparringScore,
    };
