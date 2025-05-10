// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'common.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Shot _$ShotFromJson(Map<String, dynamic> json) =>
    Shot((json['shot_ordinal'] as num).toInt(), json['score'] as String?);

Map<String, dynamic> _$ShotToJson(Shot instance) => <String, dynamic>{
  'shot_ordinal': instance.shotOrdinal,
  'score': instance.score,
};
