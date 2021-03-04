using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace NFA
{
    public class DFARule
    {
        public int State;
        public int Character;
        public int NextState;

        public DFARule(int state, int character, int nextState)
        {
            State = state;
            Character = character;
            NextState = nextState;
        }

        public bool AppliesTo(int state, int character)
        {
            return State == state && Character == character;
        }
        public int Follow()
        {
            return NextState;
        }
    }

    public class DFARulebook
    {
        List<DFARule> Rules = new List<DFARule>();

        public void AddRule(int startState, int switchCharacter, int nextState)
        {
            Rules.Add(new DFARule(startState, switchCharacter, nextState));
        }

        public HashSet<DFARule> NextStates(List<int> states, int character)
        {
            HashSet<DFARule> followRules = new HashSet<DFARule>();
            foreach (int item in states)
            {
                FollowRulesFor(item, character, followRules);
            }

            return followRules;
        }

        private void FollowRulesFor(int state, int character, HashSet<DFARule> followRules)
        {
            RulesFor(state, character, followRules);
            foreach (var rule in followRules)
            {
                rule.Follow();
            }
        }

        private void RulesFor(int state, int character, HashSet<DFARule> followRules)
        {
            foreach (var item in followRules)
            {
                if (item.AppliesTo(state, character))
                {
                    followRules.Add(item);
                }
            }
        }
        public List<int> FollowFreeMoves(List<int> nextStates)
        {
            var followedRules = NextStates(nextStates, -1);
            var moreStates = new List<int>();
            foreach (var item in followedRules)
            {
                moreStates.Add(item.NextState);
            }

            if (moreStates.Any(x => nextStates.Any(y => y == x)))
            {
                return nextStates;
            }
            else
            {
                var newList = new List<int>(nextStates);
                newList.AddRange(moreStates);
                return FollowFreeMoves(newList);
            }
        }
    }

    class NFA
    {
        public List<int> CurrentStates;
        public List<int> AcceptStates;
        public DFARulebook Rulebook;

        public NFA(List<int> startStates, List<int> acceptStates, DFARulebook rulebook)
        {
            CurrentStates = startStates;
            AcceptStates = acceptStates;
            Rulebook = rulebook;
        }

        public void AddAccepedState(int state)
        {
            AcceptStates.Add(state);
        }

        public void RemoveAcceptedState(int state)
        {
            AcceptStates.Remove(state);
        }

        public bool Accepting()
        {
            foreach (int acceptState in AcceptStates)
            {
                foreach(int currentState in CurrentStates)
                {
                    if (acceptState == currentState)
                    {
                        return true;
                    }
                }
            }
            return false;
        }

        public void ReadCharacter(int character)
        {
            HashSet<DFARule> nextStates = new HashSet<DFARule>();
            nextStates = Rulebook.NextStates(CurrentStates, character);
            foreach (DFARule rule in nextStates)
            {
                CurrentStates.Add(rule.NextState);
            }
        }

        public void ReadString(string str)
        {
            foreach (char c in str)
            {
                int character = c;
                ReadCharacter(c);
            }
        }

        public void CheckAndPrintWords(List<string> words)
        {
            List<int> start = new List<int>();
            start.Add(1);
            foreach (string word in words)
            {
                NFA newNFA = new NFA(start, AcceptStates, Rulebook);
                newNFA.ReadString(word);
                if (newNFA.Accepting())
                {
                    Console.WriteLine(word);
                }
            }
            
        }

    }
    class Program
    {
        static void Main(string[] args)
        {
            //DFARulebook rulebook = new DFARulebook();
            //rulebook.AddRule(2, 'b', 3);
            //rulebook.AddRule(2, 'a', 3);
            //rulebook.AddRule(1, 'b', 1);
            //rulebook.AddRule(1, 'a', 1);
            //rulebook.AddRule(1, 'b', 2);
            //rulebook.AddRule(3, 'a', 4);
            //rulebook.AddRule(3, 'b', 4);

            List<int> startState = new List<int>();
            startState.Add(1);
            //List<int> state = new List<int>();
            //state.Add(4);
            //NFA dfa = new NFA(startState, state, rulebook);
            //dfa.ReadString("bbabb");
            //if (dfa.Accepting())
            //{
            //    Console.WriteLine("True");
            //    Console.ReadLine();
            //}
            //else
            //{
            //    Console.WriteLine("False");
            //    Console.ReadLine();
            //}

            string alphabet = "абвгдеёжзийклмнопрстуфхцчшщъэьэюя";
            DFARulebook alfbook = new DFARulebook();
            alfbook.AddRule(1, 'a', 2);
            alfbook.AddRule(1, 'о', 2);
            alfbook.AddRule(1, 'и', 2);
            alfbook.AddRule(1, 'е', 3);
            foreach (char c in alphabet)
            {
                alfbook.AddRule(1,c,1);
            }
            List<int> stateEnd = new List<int>();
            stateEnd.Add(3);
            NFA alfNfa = new NFA(startState, stateEnd, alfbook);
            List<string> words = new List<string>();
            words.Add("большое");
            words.Add("тигр");
            words.Add("большие");
            words.Add("серый");
            words.Add("урожае");
            alfNfa.CheckAndPrintWords(words);
        }
    }
}
