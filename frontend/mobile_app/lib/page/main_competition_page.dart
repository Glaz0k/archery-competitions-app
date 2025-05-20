import 'dart:async';
import 'dart:developer';

import 'package:flutter/material.dart';
import 'package:mobile_app/api/api.dart';
import 'package:mobile_app/api/responses.dart';
import 'package:mobile_app/page/widgets/competition_field.dart';
import 'package:mobile_app/page/widgets/onion_bar.dart';
import 'package:provider/provider.dart';

const months = [
  'января',
  'февраля',
  'марта',
  'апреля',
  'мая',
  'июня',
  'июля',
  'августа',
  'сентября',
  'октября',
  'ноября',
  'декабря',
];

class MainCompetitionPage extends StatefulWidget {
  const MainCompetitionPage({super.key});

  @override
  State<MainCompetitionPage> createState() => _MainCompetitionPage();
}

class _MainCompetitionPage extends State<MainCompetitionPage> {
  final List<FullCompetitionData> _competitions = [];
  @override
  void initState() {
    super.initState();
    var api = context.read<Api>();
    _loadData(api);
  }

  Future<void> _loadData(Api api) async {
    log("Загружаем данные");
    for (var cup in await api.getCups()) {
      for (var competition in await api.getCupsCompetitions(cup.id)) {
        for (var individualGroup in await api.getCompetitionsIndividualGroups(competition.id)) {
          setState(() {
            _competitions.add(FullCompetitionData(cup, competition, individualGroup));
          });
        }
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    var api = context.watch<Api>();
    return RefreshIndicator(
      onRefresh: () {
        _competitions.clear();
        return _loadData(api);
      },
      child: Scaffold(
        appBar: OnionBar("Соревнования", context),
        body: Center(
          child: ListView.builder(
            itemCount: _competitions.length,
            itemBuilder: (context, index) {
              final data = _competitions[index];
              return CompetitionField(
                nameOfComp: data.competition.stage.name,
                date: _formatCompetitionDate(
                  data.competition.startDate,
                  data.competition.endDate,
                ),
                groupId: data.individualGroup.id
              );
            },
          ),
        ),
      ),
    );
  }

  String _formatCompetitionDate(String? startStr, String? endStr) {
    if (startStr == null && endStr == null) {
      return "даты соревнований не указаны";
    }

    DateTime? start = startStr != null ? DateTime.parse(startStr) : null;
    DateTime? end = endStr != null ? DateTime.parse(endStr) : null;

    if (start == null) return 'до ${formatDate(end!)}';
    if (end == null) return 'с ${formatDate(start)}';

    if (start.day == end.day &&
        start.month == end.month &&
        start.year == end.year) {
      return formatDate(start);
    }

    if (start.year == end.year && start.month == end.month) {
      return '${start.day}-${end.day} ${getMonth(start.month)} ${start.year} г.';
    }

    if (start.year == end.year) {
      return '${start.day} ${getMonth(start.month)} - '
          '${end.day} ${getMonth(end.month)} ${start.year} г.';
    }

    return '${formatDate(start)} - ${formatDate(end)}';
  }
}

String getMonth(int month) {
  return months[month - 1];
}

String formatDate(DateTime date) {
  return '${date.day} ${getMonth(date.month)} ${date.year} г.';
}

class FullCompetitionData {
  final Cup cup;
  final Competition competition;
  final IndividualGroup individualGroup;

  FullCompetitionData(this.cup, this.competition, this.individualGroup);

}