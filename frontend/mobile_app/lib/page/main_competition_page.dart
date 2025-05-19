import 'package:flutter/material.dart';
import 'package:mobile_app/api/api.dart';
import 'package:mobile_app/api/responses.dart';
import 'package:mobile_app/page/widgets/CompetitionField.dart';
import 'package:mobile_app/page/widgets/onion_bar.dart';
import 'package:provider/provider.dart';

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
              );
            },
          ),
        ),
      ),
    );
  }

  String _formatCompetitionDate(DateTime? start, DateTime? end) {
    if (start == null && end == null) {
      return "соревнований не указаны";
    }
    if (start == null) return 'до ${_formatDate(end!)}';
    if (end == null) return 'с ${_formatDate(start)}';

    if (start.day == end.day &&
        start.month == end.month &&
        start.year == end.year) {
      return _formatDate(start);
    }

    if (start.year == end.year && start.month == end.month) {
      return '${start.day}-${end.day} ${_getMonth(start.month)} ${start.year} г.';
    }

    if (start.year == end.year) {
      return '${start.day} ${_getMonth(start.month)} - '
          '${end.day} ${_getMonth(end.month)} ${start.year} г.';
    }

    return '${_formatDate(start)} - ${_formatDate(end)}';
  }

  String _getMonth(int month) {
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
    return months[month - 1];
  }

  String _formatDate(DateTime date) {
    return '${date.day} ${_getMonth(date.month)} ${date.year} г.';
  }
}
