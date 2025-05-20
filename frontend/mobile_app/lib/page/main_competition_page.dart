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
  Map<int, List<Competition>> _competitions = {};
  List<Cup> _cups = [];
  Map<int, List<IndividualGroup>> _individualGroups = {};

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    final api = context.read<Api>();
    final cups = await api.getCups();
    if (cups.isEmpty) return;

    final competitionsFuture = cups.map(
      (cup) => api.getCupsCompetitions(cup.id),
    );
    final allCompetitions = await Future.wait(competitionsFuture);

    final cupCompetitions = {
      for (int i = 0; i < cups.length; i++) cups[i].id: allCompetitions[i],
    };

    final allGroups = await Future.wait(
      allCompetitions
          .expand((competitions) => competitions)
          .map(
            (competition) =>
                api.getCompetitionsIndividualGroups(competition.id),
          ),
    );
    final individualGroups = {
      for (int i = 0; i < allCompetitions.expand((c) => c).length; i++)
        allCompetitions.expand((c) => c).elementAt(i).id: allGroups[i],
    };
    setState(() {
      _cups = cups;
      _competitions = cupCompetitions;
      _individualGroups = individualGroups;
    });
  }

  @override
  Widget build(BuildContext context) {
    final allCompetitions =
        _competitions.values.expand((list) => list).toList();
    return RefreshIndicator(
      onRefresh: _loadData,
      child: Scaffold(
        appBar: OnionBar("Соревнования", context),
        body: Center(
          child: ListView.builder(
            itemCount: allCompetitions.length,
            itemBuilder: (context, index) {
              final competition = allCompetitions[index];
              return CompetitionField(
                nameOfComp: competition.stage.name,
                date: _formatCompetitionDate(
                  competition.startDate,
                  competition.endDate,
                ),
                // Todo: Как-нибудь достань мне этот айдишник. Этот способ не работает, выкидывает NullPointerException
                groupId: _individualGroups[index]!.first.id
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
